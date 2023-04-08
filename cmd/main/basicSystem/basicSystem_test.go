package basicsystem

import (
	manager "OCSS/FoosballGeneticLearning/pkg/Manager"
	"os"
	"testing"
)

func TestBasicSystem(t *testing.T) {
	targetSystem := BasicSystem{}
	manager := manager.NewManager(&targetSystem, 100, false)
	for generationIndex := 0; generationIndex < 100; generationIndex++ {
		manager.SimulateGeneration()
	}

	os.RemoveAll("data")
	os.RemoveAll("logs")
}
