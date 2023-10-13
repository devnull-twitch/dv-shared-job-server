package sharedjob

import (
	"slices"
	"sync"

	"github.com/sirupsen/logrus"
)

type (
	JobType string
	Job     struct {
		ID                  string    `json:"id"`
		JobType             JobType   `json:"type"`
		StartingStationName StationID `json:"-"`
		StartingTrack       string    `json:"starting_track"`
		startTrackType      TrackTypeID
		TargetStationName   StationID `json:"target_station"`
		TargetTrack         string    `json:"target_track"`
		targetTrackType     TrackTypeID
		CarCount            int       `json:"car_count"`
		CargoType           CargoType `json:"cargo_type"`
		Wage                int       `json:"wage"`
		jobReserved         bool
		jobActive           bool
		jobAssignedUser     string
		jobSpawned          bool
	}
)

func (j *Job) IsReserved() bool {
	return j.jobReserved
}

func (j *Job) IsActive() bool {
	return j.jobActive
}

func (j *Job) IsSpawned() bool {
	return j.jobSpawned
}

func (j *Job) IsAssigned() bool {
	return j.jobAssignedUser != ""
}

func (j *Job) GetAssignedUser() string {
	return j.jobAssignedUser
}

func (j *Job) GetStartStation() *LogicStation {
	return GetStation(j.StartingStationName)
}

const (
	LogisticHaulJobType   JobType = "logistics"
	ShuntingLoadJobType   JobType = "shunting_load"
	ShuntingUnloadJobType JobType = "shunting_unload"
	FreightJobType        JobType = "freight"
)

func (t JobType) AsID() string {
	switch t {
	case LogisticHaulJobType:
		return "SLH"
	case ShuntingLoadJobType:
		return "SSL"
	case ShuntingUnloadJobType:
		return "SSU"
	case FreightJobType:
		return "SFH"
	}

	return "UNK"
}

func (t JobType) String() string {
	return string(t)
}

var jobLock = sync.Mutex{}

/**
 * - Do we allow conflicting target tracks on available jobs. Just not on activbe jobs
 *   - Con: A lot more action despawning and spawning jobs
 *   - Pro: More jobs overall available
 * - Dynamic car list based on station needs.
 * - Stations queue jobs based on input. If tracks are available spawn job. If a job despawns we need to queu it again.
 *   Or maybe we never take it out of quue as long as no user accepted it. More checks for spawing maybe ... seems like
 *   the better idea.
 */

func GetAllStationJobs(sourceStation StationID) []Job {
	stationJobs := make([]Job, 0)
	logicStation := GetStation(sourceStation)
	for _, j := range logicStation.JobQueue {
		if j.jobSpawned {
			stationJobs = append(stationJobs, *j)
		}
	}

	return stationJobs
}

func GetAllStationJobsForUsername(sourceStation StationID, userName string) []Job {
	stationJobs := make([]Job, 0)
	logicStation := GetStation(sourceStation)
	for _, j := range logicStation.JobQueue {
		if j.jobSpawned && (!j.jobActive || j.jobAssignedUser == userName) {
			stationJobs = append(stationJobs, *j)
		}
	}

	return stationJobs
}

func ReserveJob(userName string, jobID string) bool {
	for _, logicStation := range AllStations {
		for index, j := range logicStation.JobQueue {
			if j.ID == jobID {
				if j.jobReserved || j.jobActive {
					return false
				}

				jobLock.Lock()
				defer jobLock.Unlock()

				logicStation.JobQueue[index].jobReserved = true
				logicStation.JobQueue[index].jobAssignedUser = userName

				return true
			}
		}
	}

	return false
}

func TakeJob(userName string, jobID string, progressCh chan<- ProgressMessage) (bool, []*Job, []*Job, []*Job) {
	jobLock.Lock()
	defer func() {
		jobLock.Unlock()
	}()

	for _, logicStation := range AllStations {
		for index, j := range logicStation.JobQueue {
			if j.ID == jobID {
				if !j.jobReserved || j.jobActive || j.jobAssignedUser != userName {
					return false, nil, nil, nil
				}

				logicStation.JobQueue[index].jobActive = true

				updateData := updateAllJobs(j)
				if !slices.ContainsFunc(updateData.changedJobs, func(checkJob *Job) bool {
					return checkJob.ID == jobID
				}) {
					updateData.changedJobs = slices.Insert(updateData.changedJobs, 0, j)
				}

				for _, sid := range updateData.notifyStationIDs {
					logrus.WithField("station_id", sid).Info("notifying station after job completion impact")
					progressCh <- ProgressMessage{StationID: sid}
				}

				return true, updateData.unspawnJobs, updateData.spawnJobs, updateData.changedJobs
			}
		}
	}

	jobLock.Unlock()
	return false, nil, nil, nil
}

func FinishJob(userName string, jobID string, progressCh chan<- ProgressMessage) (bool, []*Job, []*Job, []*Job, []*Job) {
	for _, logicStation := range AllStations {
		for index, j := range logicStation.JobQueue {
			if j.ID == jobID {
				if !j.jobActive || j.jobAssignedUser != userName {
					return false, nil, nil, nil, nil
				}

				jobLock.Lock()
				defer func() {
					jobLock.Unlock()
				}()

				filtered := make([]*Job, 0, len(logicStation.JobQueue))
				filtered = append(filtered, logicStation.JobQueue[:index]...)
				filtered = append(filtered, logicStation.JobQueue[index+1:]...)
				logicStation.JobQueue = filtered

				newlyCreatedJobs := GetStation(j.TargetStationName).ProcessJob(j)

				updateData := updateAllJobs(j)

				// fot now unspawn the finished job and spawn follow up as fresh job.
				// would be sweet if we would keep the cars
				// ( aka make use of proper jobChains in game logic, but hey ...  )
				updateData.unspawnJobs = slices.Insert(updateData.unspawnJobs, 0, j)

				for _, sid := range updateData.notifyStationIDs {
					logrus.WithField("station_id", sid).Info("notifying station after job completion impact")
					progressCh <- ProgressMessage{StationID: sid}
				}

				return true, updateData.unspawnJobs, updateData.spawnJobs, updateData.changedJobs, newlyCreatedJobs
			}
		}
	}

	return false, nil, nil, nil, nil
}

type updateAllPayload struct {
	notifyStationIDs []StationID
	unspawnJobs      []*Job
	spawnJobs        []*Job
	changedJobs      []*Job
}

func updateAllJobs(srcJob *Job) updateAllPayload {
	retVal := updateAllPayload{
		notifyStationIDs: make([]StationID, 0),
		unspawnJobs:      make([]*Job, 0),
		spawnJobs:        make([]*Job, 0),
		changedJobs:      make([]*Job, 0),
	}

	changeFlagPerStation := make(map[StationID]bool)
	for stationID, logicStation := range AllStations {
		changed, stationDeleteIds, stationNewJobs, stationChangedJobs := logicStation.ValidateJobs(ShuntingUnloadJobType)
		changeFlagPerStation[stationID] = changed

		retVal.unspawnJobs = append(retVal.unspawnJobs, stationDeleteIds...)
		retVal.spawnJobs = append(retVal.spawnJobs, stationNewJobs...)
		retVal.changedJobs = append(retVal.changedJobs, stationChangedJobs...)
	}
	for stationID, logicStation := range AllStations {
		changed, stationDeleteIds, stationNewJobs, stationChangedJobs := logicStation.ValidateJobs(FreightJobType)
		changeFlagPerStation[stationID] = changed || changeFlagPerStation[stationID]

		retVal.unspawnJobs = append(retVal.unspawnJobs, stationDeleteIds...)
		retVal.spawnJobs = append(retVal.spawnJobs, stationNewJobs...)
		retVal.changedJobs = append(retVal.changedJobs, stationChangedJobs...)
	}
	for stationID, logicStation := range AllStations {
		changed, stationDeleteIds, stationNewJobs, stationChangedJobs := logicStation.ValidateJobs(ShuntingLoadJobType)
		changeFlagPerStation[stationID] = changed || changeFlagPerStation[stationID]

		retVal.unspawnJobs = append(retVal.unspawnJobs, stationDeleteIds...)
		retVal.spawnJobs = append(retVal.spawnJobs, stationNewJobs...)
		retVal.changedJobs = append(retVal.changedJobs, stationChangedJobs...)
	}
	for stationID, logicStation := range AllStations {
		changed, stationDeleteIds, stationNewJobs, stationChangedJobs := logicStation.ValidateJobs(LogisticHaulJobType)
		changeFlagPerStation[stationID] = changed || changeFlagPerStation[stationID]

		retVal.unspawnJobs = append(retVal.unspawnJobs, stationDeleteIds...)
		retVal.spawnJobs = append(retVal.spawnJobs, stationNewJobs...)
		retVal.changedJobs = append(retVal.changedJobs, stationChangedJobs...)
	}

	for stationID, isChanged := range changeFlagPerStation {
		if stationID == srcJob.StartingStationName || stationID == srcJob.TargetStationName {
			retVal.notifyStationIDs = append(retVal.notifyStationIDs, stationID)
			continue
		}

		if !isChanged {
			continue
		}

		for _, stationJob := range AllStations[stationID].JobQueue {
			if stationJob.TargetStationName == srcJob.TargetStationName {
				retVal.notifyStationIDs = append(retVal.notifyStationIDs, stationID)
			}
		}
	}

	retVal.notifyStationIDs = slices.Compact(retVal.notifyStationIDs)

	return retVal
}
