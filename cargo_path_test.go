package sharedjob

import (
	"fmt"
	"testing"
)

func TestCargoValidation(t *testing.T) {
	// setup station processors
	Setup()

	// make sure every input is accounted for
	for _, station := range AllStations {
		for _, proc := range station.Processor {
			for _, input := range proc.allowedInput {
				if input == None {
					continue
				}

				if !findOutputTo(station, input) {
					t.Error(fmt.Errorf("no output for this %s from %s", string(station.ID), string(input)))
				}
			}
		}
	}
}

func findOutputTo(target *LogicStation, cargo CargoType) bool {
	for _, station := range AllStations {
		for _, proc := range station.Processor {
			if proc.output == cargo {
				for _, targetStation := range proc.targetStations {
					if targetStation == target.ID {
						return true
					}
				}
			}
		}
	}

	return false
}
