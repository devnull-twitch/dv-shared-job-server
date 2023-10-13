package sharedjob

import (
	"fmt"
	"math/rand"
)

type (
	StationID   string
	YardID      string
	TrackNumber string
	TrackTypeID string

	YardTracks    map[YardID][]TrackNumber
	TrackYards    map[TrackTypeID]YardTracks
	StationTracks map[StationID]TrackYards
)

const (
	OutputTrackType  TrackTypeID = "O"
	StorageTrackType TrackTypeID = "S"
	InputTrackType   TrackTypeID = "I"
	LoadingTrackType TrackTypeID = "L"

	StationCSW StationID = "CSW"
	StationCM  StationID = "CM"
	StationFF  StationID = "FF"
	StationFM  StationID = "FM"
	StationFRC StationID = "FRC"
	StationFRS StationID = "FRS"
	StationGF  StationID = "GF"
	StationHB  StationID = "HB"
	StationHMB StationID = "HMB"
	StationIME StationID = "IME"
	StationIMW StationID = "IMW"
	StationMF  StationID = "MF"
	StationMB  StationID = "MB"
	StationOWC StationID = "OWC"
	StationOWN StationID = "OWN"
	StationSW  StationID = "SW"
	StationSM  StationID = "SM"
)

func (s StationID) String() string {
	return string(s)
}

var StationTrackData = StationTracks{
	StationCSW: TrackYards{
		InputTrackType: YardTracks{
			YardID("C"): []TrackNumber{"02", "03"},
		},
		StorageTrackType: YardTracks{
			YardID("C"): []TrackNumber{"04"},
		},
		OutputTrackType: YardTracks{
			YardID("C"): []TrackNumber{"05"},
		},
		LoadingTrackType: YardTracks{
			YardID("C"): []TrackNumber{"06"},
		},
	},
	StationCM: TrackYards{
		InputTrackType: YardTracks{
			YardID("B"): []TrackNumber{"01"},
		},
		StorageTrackType: YardTracks{
			YardID("B"): []TrackNumber{"05"},
			YardID("C"): []TrackNumber{"01", "03"},
		},
		OutputTrackType: YardTracks{
			YardID("B"): []TrackNumber{"02", "03"},
		},
		LoadingTrackType: YardTracks{
			YardID("A"): []TrackNumber{"03"},
		},
	},
	StationFF: TrackYards{
		InputTrackType: YardTracks{
			YardID("C"): []TrackNumber{"04", "06"},
			YardID("D"): []TrackNumber{"02"},
		},
		StorageTrackType: YardTracks{
			YardID("A"): []TrackNumber{"01"},
			YardID("C"): []TrackNumber{"01"},
			YardID("D"): []TrackNumber{"03", "04"},
		},
		OutputTrackType: YardTracks{
			YardID("C"): []TrackNumber{"02", "03", "05", "07", "08"},
		},
		LoadingTrackType: YardTracks{
			YardID("D"): []TrackNumber{"01"},
		},
	},
	StationFM: TrackYards{
		InputTrackType: YardTracks{
			YardID("B"): []TrackNumber{"02"},
		},
		StorageTrackType: YardTracks{
			YardID("B"): []TrackNumber{"01", "03"},
		},
		OutputTrackType: YardTracks{
			YardID("B"): []TrackNumber{"05", "06"},
		},
		LoadingTrackType: YardTracks{
			YardID("A"): []TrackNumber{"01", "02", "03"},
		},
	},
	StationFRC: TrackYards{
		StorageTrackType: YardTracks{
			YardID("B"): []TrackNumber{"04"},
			YardID("C"): []TrackNumber{"01", "02"},
		},
		OutputTrackType: YardTracks{
			YardID("B"): []TrackNumber{"02"},
			YardID("C"): []TrackNumber{"04"},
		},
		LoadingTrackType: YardTracks{
			YardID("B"): []TrackNumber{"01"},
		},
		InputTrackType: YardTracks{},
	},
	StationGF: TrackYards{
		InputTrackType: YardTracks{
			YardID("D"): []TrackNumber{"05", "06"},
		},
		StorageTrackType: YardTracks{
			YardID("A"): []TrackNumber{"02", "03"},
			YardID("B"): []TrackNumber{"02", "03"},
			YardID("D"): []TrackNumber{"01"},
		},
		OutputTrackType: YardTracks{
			YardID("D"): []TrackNumber{"02", "03", "04"},
		},
		LoadingTrackType: YardTracks{
			YardID("B"): []TrackNumber{"01"},
		},
	},
	StationHB: TrackYards{
		InputTrackType: YardTracks{
			YardID("C"): []TrackNumber{"02"},
			YardID("D"): []TrackNumber{"04"},
			YardID("E"): []TrackNumber{"08", "09"},
			YardID("G"): []TrackNumber{"05"},
		},
		StorageTrackType: YardTracks{
			YardID("C"): []TrackNumber{"01"},
			YardID("D"): []TrackNumber{"01", "02", "05"},
			YardID("G"): []TrackNumber{"01", "02", "06", "07"},
		},
		OutputTrackType: YardTracks{
			YardID("D"): []TrackNumber{"03", "06"},
			YardID("E"): []TrackNumber{"01", "02", "03", "04", "05", "07", "10", "11"},
			YardID("G"): []TrackNumber{"03"},
		},
		LoadingTrackType: YardTracks{
			YardID("C"): []TrackNumber{"03"},
			YardID("D"): []TrackNumber{"07"},
		},
	},
	StationIME: TrackYards{
		InputTrackType: YardTracks{
			YardID("C"): []TrackNumber{"04"},
		},
		StorageTrackType: YardTracks{
			YardID("B"): []TrackNumber{"01"},
			YardID("C"): []TrackNumber{"01"},
		},
		OutputTrackType: YardTracks{
			YardID("B"): []TrackNumber{"02", "04"},
			YardID("C"): []TrackNumber{"03"},
		},
		LoadingTrackType: YardTracks{
			YardID("A"): []TrackNumber{"01"},
		},
	},
	StationIMW: TrackYards{
		InputTrackType: YardTracks{
			YardID("B"): []TrackNumber{"02"},
		},
		StorageTrackType: YardTracks{
			YardID("B"): []TrackNumber{"01", "07"},
		},
		OutputTrackType: YardTracks{
			YardID("B"): []TrackNumber{"03", "04", "06"},
		},
		LoadingTrackType: YardTracks{
			YardID("B"): []TrackNumber{"08"},
		},
	},
	StationMF: TrackYards{
		InputTrackType: YardTracks{
			YardID("C"): []TrackNumber{"03", "04"},
		},
		StorageTrackType: YardTracks{
			YardID("B"): []TrackNumber{"01", "06"},
			YardID("C"): []TrackNumber{"02"},
		},
		OutputTrackType: YardTracks{
			YardID("B"): []TrackNumber{"02", "04", "05"},
		},
		LoadingTrackType: YardTracks{
			YardID("C"): []TrackNumber{"01"},
		},
	},
	StationMB: TrackYards{
		InputTrackType:   YardTracks{},
		StorageTrackType: YardTracks{},
		OutputTrackType:  YardTracks{},
		LoadingTrackType: YardTracks{},
	},
	StationOWC: TrackYards{
		InputTrackType: YardTracks{},
		StorageTrackType: YardTracks{
			YardID("A"): []TrackNumber{"02", "03"},
			YardID("B"): []TrackNumber{"06"},
		},
		OutputTrackType: YardTracks{
			YardID("B"): []TrackNumber{"01", "03", "04", "05"},
		},
		LoadingTrackType: YardTracks{
			YardID("A"): []TrackNumber{"01"},
		},
	},
	StationOWN: TrackYards{
		InputTrackType: YardTracks{},
		StorageTrackType: YardTracks{
			YardID("B"): []TrackNumber{"02"},
			YardID("C"): []TrackNumber{"01"},
		},
		OutputTrackType: YardTracks{
			YardID("B"): []TrackNumber{"03", "04", "05"},
			YardID("C"): []TrackNumber{"03"},
		},
		LoadingTrackType: YardTracks{
			YardID("B"): []TrackNumber{"06"},
		},
	},
	StationSW: TrackYards{
		StorageTrackType: YardTracks{
			YardID("B"): []TrackNumber{"01"},
			YardID("C"): []TrackNumber{"04"},
		},
		OutputTrackType: YardTracks{
			YardID("C"): []TrackNumber{"01"},
		},
		InputTrackType: YardTracks{
			YardID("C"): []TrackNumber{"03"},
			YardID("B"): []TrackNumber{"03"},
		},
		LoadingTrackType: YardTracks{
			YardID("B"): []TrackNumber{"04"},
		},
	},
	StationSM: TrackYards{
		LoadingTrackType: YardTracks{
			YardID("A"): []TrackNumber{"07"},
		},
		StorageTrackType: YardTracks{
			// YardID("A"): []TrackNumber{"03", "04", "05"},
			YardID("B"): []TrackNumber{"07", "08"},
		},
		InputTrackType: YardTracks{
			YardID("A"): []TrackNumber{"06"},
			YardID("B"): []TrackNumber{"03"},
		},
		OutputTrackType: YardTracks{
			YardID("B"): []TrackNumber{"01", "02", "04", "06"},
		},
	},
}

func (s *LogicStation) GetAllFullTrackNames(yardType TrackTypeID) []string {
	fullNames := make([]string, 0)
	yardTracks := StationTrackData[s.ID][yardType]
	for _, yardId := range []YardID{"A", "B", "C", "D", "E", "F", "G", "H"} {
		numbers, ok := yardTracks[yardId]
		if !ok {
			continue
		}

		for _, numStr := range numbers {
			fullNames = append(fullNames, fmt.Sprintf("%s-%s-%s-%s", s.ID, yardId, numStr, yardType))
		}
	}

	return fullNames
}

func (s *LogicStation) GetFreeTrackName(yardType TrackTypeID) *string {
	trackNames := s.GetAllFullTrackNames(yardType)

	for _, trackName := range trackNames {
		if s.isTrackFree(trackName) {
			return &trackName
		}
	}

	return nil
}

func (s *LogicStation) GetRandomTrackName(yardType TrackTypeID) string {
	trackNames := s.GetAllFullTrackNames(yardType)
	if len(trackNames) == 0 {
		return ""
	}

	index := rand.Intn(len(trackNames))
	return trackNames[index]
}

func (s *LogicStation) isTrackFree(trackName string) bool {
	for _, station := range AllStations {
		for _, j := range station.JobQueue {
			if j.jobSpawned && j.StartingTrack == trackName {
				return false
			}

			if j.jobActive && j.TargetTrack == trackName {
				return false
			}
		}
	}

	return true
}
