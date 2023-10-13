package sharedjob

import (
	"fmt"
	"math/rand"
	"slices"

	"github.com/sirupsen/logrus"
)

const (
	MAX_CARS_PER_JOB int = 12
)

type StationProcessor struct {
	allowedInput   []CargoType
	blueprint      map[CargoType]int
	buffer         map[CargoType]int
	output         CargoType
	targetStations []StationID
}

type LogicStation struct {
	ID            StationID
	JobQueue      []*Job
	lastJobNum    int
	lastProcIndex int
	Processor     []*StationProcessor
	cargoBuffer   map[CargoType]int
}

var AllStations = map[StationID]*LogicStation{
	StationCSW: NewStation(StationCSW),
	StationCM:  NewStation(StationCM),
	StationFF:  NewStation(StationFF),
	StationFM:  NewStation(StationFM),
	StationFRC: NewStation(StationFRC),
	StationFRS: NewStation(StationFRS),
	StationGF:  NewStation(StationGF),
	StationHB:  NewStation(StationHB),
	StationHMB: NewStation(StationHMB),
	StationIME: NewStation(StationIME),
	StationIMW: NewStation(StationIMW),
	StationMF:  NewStation(StationMF),
	StationMB:  NewStation(StationMB),
	StationOWC: NewStation(StationOWC),
	StationOWN: NewStation(StationOWN),
	StationSW:  NewStation(StationSW),
	StationSM:  NewStation(StationSM),
}

func Setup() {
	// CSW
	// Generative output
	AllStations[StationCSW].AddProcessor(NewProcessor(map[CargoType]int{None: 1}, ScrapMetal, StationSM))

	// CM
	// Generative output
	AllStations[StationCM].AddProcessor(NewProcessor(map[CargoType]int{None: 1}, Coal, StationSM))

	// FF
	// Transformative processors
	AllStations[StationFF].AddProcessor(NewProcessor(map[CargoType]int{Wheat: 1}, Alcohol, StationHB))
	AllStations[StationFF].AddProcessor(NewProcessor(map[CargoType]int{Pigs: 1}, CannedFood, StationHB, StationCSW))
	AllStations[StationFF].AddProcessor(NewProcessor(map[CargoType]int{Chickens: 1}, CatFood, StationHB, StationCSW))

	// FM
	// Generative output
	AllStations[StationFM].AddProcessor(NewProcessor(map[CargoType]int{None: 1}, Pigs, StationFF))
	AllStations[StationFM].AddProcessor(NewProcessor(map[CargoType]int{None: 1}, Sheep, StationFF))
	AllStations[StationFM].AddProcessor(NewProcessor(map[CargoType]int{None: 1}, Chickens, StationFF))
	AllStations[StationFM].AddProcessor(NewProcessor(map[CargoType]int{None: 1}, Cows, StationFF))
	AllStations[StationFM].AddProcessor(NewProcessor(map[CargoType]int{None: 1}, Wheat, StationFF))

	// FRC
	// Generative output
	AllStations[StationFRC].AddProcessor(NewProcessor(map[CargoType]int{None: 1}, Logs, StationSW))

	// FRS
	// Generative output
	AllStations[StationFRS].AddProcessor(NewProcessor(map[CargoType]int{None: 1}, Logs, StationSW))

	// GF
	// Generative output
	AllStations[StationGF].AddProcessor(NewProcessor(map[CargoType]int{SteelBillets: 1}, ToolsIskar, StationMF, StationCSW))
	AllStations[StationGF].AddProcessor(NewProcessor(map[CargoType]int{SteelBillets: 1}, ToolsBrohm, StationMF, StationCSW))
	AllStations[StationGF].AddProcessor(NewProcessor(map[CargoType]int{SteelBillets: 1}, ToolsAAG, StationMF, StationCSW))
	AllStations[StationGF].AddProcessor(NewProcessor(map[CargoType]int{SteelBillets: 1}, ToolsNovae, StationMF, StationCSW))
	AllStations[StationGF].AddProcessor(NewProcessor(map[CargoType]int{SteelBillets: 1}, ToolsTraeg, StationMF, StationCSW))

	// HB
	// Generative output
	AllStations[StationHB].AddProcessor(NewProcessor(map[CargoType]int{None: 1}, Ammonia, StationFF))
	AllStations[StationHB].AddProcessor(NewProcessor(map[CargoType]int{None: 1}, ImportedNewCars, StationCSW))
	AllStations[StationHB].AddProcessor(NewProcessor(map[CargoType]int{None: 1}, ClothingNeoGamma, StationCSW))
	AllStations[StationHB].AddProcessor(NewProcessor(map[CargoType]int{None: 1}, Medicine, StationCSW))
	AllStations[StationHB].AddProcessor(NewProcessor(map[CargoType]int{None: 1}, ClothingNovae, StationCSW))
	AllStations[StationHB].AddProcessor(NewProcessor(map[CargoType]int{None: 1}, Acetylene, StationGF))
	AllStations[StationHB].AddProcessor(NewProcessor(map[CargoType]int{None: 1}, CryoHydrogen, StationGF))
	AllStations[StationHB].AddProcessor(NewProcessor(map[CargoType]int{None: 1}, CryoOxygen, StationGF))
	AllStations[StationHB].AddProcessor(NewProcessor(map[CargoType]int{None: 1}, Methane, StationGF))

	// IME
	// Generative output
	AllStations[StationIME].AddProcessor(NewProcessor(map[CargoType]int{None: 1}, IronOre, StationSM))

	// IMW
	// Generative output
	AllStations[StationIMW].AddProcessor(NewProcessor(map[CargoType]int{None: 1}, IronOre, StationSM))

	// SW
	// Generative output
	AllStations[StationSW].AddProcessor(NewProcessor(map[CargoType]int{Logs: 1}, Boards, StationGF))
	AllStations[StationSW].AddProcessor(NewProcessor(map[CargoType]int{Logs: 1}, Plywood, StationGF))

	// SM
	// Transformative output
	AllStations[StationSM].AddProcessor(NewProcessor(map[CargoType]int{Coal: 1, IronOre: 2}, SteelSlabs, StationGF))
	AllStations[StationSM].AddProcessor(NewProcessor(map[CargoType]int{Coal: 1, IronOre: 2}, SteelBillets, StationGF))

	// make sure we always spawn the later end of the job chain before any earlier stages
	for _, logicStation := range AllStations {
		logicStation.ValidateJobs(ShuntingUnloadJobType)
	}
	for _, logicStation := range AllStations {
		logicStation.ValidateJobs(FreightJobType)
	}
	for _, logicStation := range AllStations {
		logicStation.ValidateJobs(ShuntingLoadJobType)
	}
	for _, logicStation := range AllStations {
		logicStation.ValidateJobs(LogisticHaulJobType)
	}
}

func NewProcessor(in map[CargoType]int, out CargoType, targetStation ...StationID) *StationProcessor {
	proc := &StationProcessor{
		output:         out,
		blueprint:      in,
		allowedInput:   make([]CargoType, 0, len(in)),
		buffer:         make(map[CargoType]int),
		targetStations: targetStation,
	}
	for cType := range in {
		proc.allowedInput = append(proc.allowedInput, cType)
	}

	return proc
}

func NewStation(id StationID) *LogicStation {
	lStation := &LogicStation{
		ID:          id,
		JobQueue:    []*Job{},
		cargoBuffer: make(map[CargoType]int),
		Processor:   []*StationProcessor{},
	}

	return lStation
}

func (s *LogicStation) AddProcessor(proc *StationProcessor) {
	s.Processor = append(s.Processor, proc)

	if len(proc.allowedInput) <= 0 || (len(proc.allowedInput) == 1 && proc.allowedInput[0] == None) {
		// TODO: evil hardcoded magical numbers. make them go away!
		s.AddJob(s.ID, ShuntingLoadJobType, 5, proc.output, 9000)
	}
}

func GetStation(id StationID) *LogicStation {
	return AllStations[id]
}

func (s *LogicStation) GetJob(jobID string) (*Job, error) {
	for _, j := range s.JobQueue {
		if j.ID == jobID {
			return j, nil
		}
	}

	return nil, fmt.Errorf("job %s not found", jobID)
}

func (s *LogicStation) AddJob(
	targetStation StationID,
	jobType JobType,
	carCount int,
	cargoType CargoType,
	wage int,
) *Job {
	var (
		startTrackType  TrackTypeID
		targetTrackType TrackTypeID
	)

	switch jobType {
	case LogisticHaulJobType:
		startTrackType = StorageTrackType
		targetTrackType = StorageTrackType
	case FreightJobType:
		startTrackType = OutputTrackType
		targetTrackType = InputTrackType
	case ShuntingLoadJobType:
		startTrackType = StorageTrackType
		targetTrackType = OutputTrackType
	case ShuntingUnloadJobType:
		startTrackType = InputTrackType
		targetTrackType = StorageTrackType
	}

	j := &Job{
		ID:                  s.GetJobID(jobType),
		StartingStationName: s.ID,
		TargetStationName:   targetStation,
		JobType:             jobType,
		CarCount:            carCount,
		CargoType:           cargoType,
		Wage:                wage,
		startTrackType:      startTrackType,
		targetTrackType:     targetTrackType,
	}

	logrus.WithFields(logrus.Fields{
		"job_id":     j.ID,
		"cargo_type": j.CargoType,
	}).Info("adding job")

	s.JobQueue = append(s.JobQueue, j)
	return j
}

func (s *LogicStation) ProcessJob(j *Job) []*Job {
	switch j.JobType {
	case ShuntingUnloadJobType:
		newJob := s.procShuntingUnload(j)
		if newJob != nil {
			return []*Job{newJob}
		}
	case ShuntingLoadJobType:
		return s.procShuntingLoad(j)
	case FreightJobType:
		newJob := s.procFreight(j)
		if newJob != nil {
			return []*Job{newJob}
		}
	}

	return []*Job{}
}

func (s *LogicStation) procFreight(j *Job) *Job {
	return s.AddJob(s.ID, ShuntingUnloadJobType, j.CarCount, j.CargoType, j.CarCount*500)
}

func (s *LogicStation) procShuntingUnload(j *Job) *Job {
	for i := 0; i < len(s.Processor); i++ {
		index := (s.lastProcIndex + i) % len(s.Processor)
		proc := s.Processor[index]

		if proc.isAllowed(j.CargoType) {
			if proc.output == None {
				return nil
			}

			proc.buffer[j.CargoType] += j.CarCount

			outCargo := proc.makeOutput()
			if outCargo > 0 {
				return s.addCargo(proc.output, outCargo)
			}
			return nil
		}
	}

	s.lastProcIndex++
	return nil
}

func (s *LogicStation) procShuntingLoad(j *Job) []*Job {
	newJobs := make([]*Job, 0)
	for _, proc := range s.Processor {
		if proc.output == j.CargoType {
			targetStation := proc.targetStations[rand.Intn(len(proc.targetStations))]
			newJobs = append(newJobs, s.AddJob(targetStation, FreightJobType, j.CarCount, j.CargoType, j.Wage))

			if len(proc.allowedInput) <= 0 || (len(proc.allowedInput) == 1 && proc.allowedInput[0] == None) {
				newJobs = append(newJobs, s.AddJob(s.ID, ShuntingLoadJobType, j.CarCount, proc.output, j.Wage))
			}

			break
		}
	}

	return newJobs
}

func (s *LogicStation) addCargo(cType CargoType, count int) *Job {
	for _, j := range s.JobQueue {
		if j.CargoType == cType {
			j.CarCount += count
			if j.CarCount > MAX_CARS_PER_JOB {
				count = j.CarCount % MAX_CARS_PER_JOB
				j.CarCount = MAX_CARS_PER_JOB
			}

			if count > 0 {
				if _, ok := s.cargoBuffer[cType]; !ok {
					s.cargoBuffer[cType] = 0
				}

				s.cargoBuffer[cType] += count
			}

			return nil
		}
	}

	return s.AddJob(s.ID, ShuntingLoadJobType, count, cType, count*1000)
}

func (s *LogicStation) trySpawnJob(j *Job) (changed, newSpawn, despawn bool) {
	defer func() {
		logrus.WithFields(logrus.Fields{
			"station": s.ID,
			"job":     j.ID,
		}).WithField("change", changed).Info("job validation")
	}()

	if j.jobActive {
		return
	}

	var startTrackPtr *string
	if !j.jobSpawned {
		startTrackPtr = s.GetFreeTrackName(j.startTrackType)
		if startTrackPtr == nil {
			return
		}
	}

	targetLogicStation := GetStation(j.TargetStationName)
	targetTrackPtr := targetLogicStation.GetFreeTrackName(j.targetTrackType)
	if targetTrackPtr == nil {
		if j.jobSpawned {
			j.jobSpawned = false
			changed = true
			despawn = true
			return
		}

		return
	}

	if j.TargetTrack != *targetTrackPtr {
		changed = true
	}
	j.TargetTrack = *targetTrackPtr

	if !j.jobSpawned {
		j.StartingTrack = *startTrackPtr
		j.jobSpawned = true
		changed = true
		newSpawn = true
		return
	}

	return
}

func (s *LogicStation) ValidateJobs(jobType JobType) (anyChanges bool, unspawnJobs []*Job, spawnJobs []*Job, changedJobs []*Job) {
	unspawnJobs = make([]*Job, 0)
	spawnJobs = make([]*Job, 0)
	changedJobs = make([]*Job, 0)

	for _, j := range s.JobQueue {
		if j.JobType != jobType {
			continue
		}

		changed, newSpawn, deSpawn := s.trySpawnJob(j)

		if !anyChanges && changed {
			anyChanges = true
		}

		if !changed {
			continue
		}

		if deSpawn {
			unspawnJobs = append(unspawnJobs, j)
			continue
		}

		if newSpawn {
			spawnJobs = append(spawnJobs, j)
			continue
		}

		changedJobs = append(changedJobs, j)
	}

	return
}

func (s *LogicStation) GetJobID(t JobType) string {
	s.lastJobNum++
	return fmt.Sprintf("%s-%s-%0d", s.ID, t.AsID(), s.lastJobNum)
}

func (s *LogicStation) AllInputs() []string {
	str := []string{}
	for _, proc := range s.Processor {
		for _, cType := range proc.allowedInput {
			if cType != None && !containsString(str, cType.String()) {
				str = append(str, cType.String())
			}
		}
	}

	return slices.Compact(str)
}

func (s *LogicStation) AllOutputs() []string {
	str := []string{}
	for _, proc := range s.Processor {
		if proc.output != None && !containsString(str, proc.output.String()) {
			str = append(str, proc.output.String())
		}
	}

	return slices.Compact(str)
}

func containsString(haystack []string, needle string) bool {
	for _, str := range haystack {
		if str == needle {
			return true
		}
	}

	return false
}

func (proc *StationProcessor) isAllowed(cType CargoType) bool {
	for _, bpCargoType := range proc.allowedInput {
		if bpCargoType == cType {
			return true
		}
	}

	return false
}

func (proc *StationProcessor) makeOutput() int {
	goon := true
	outCounter := 0
	for goon {
		for _, cargoType := range proc.allowedInput {
			if proc.buffer[cargoType] < proc.blueprint[cargoType] {
				goon = false
				break
			}
		}

		if !goon {
			break
		}

		for _, cargoType := range proc.allowedInput {
			proc.buffer[cargoType] -= proc.blueprint[cargoType]
		}

		outCounter++
	}

	return outCounter
}
