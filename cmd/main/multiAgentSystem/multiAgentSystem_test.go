package multiagentsystem

import (
	manager "OCSS/FoosballGeneticLearning/pkg/Manager"
	"os"
	"testing"
)

func TestMultiAgentSystem(t *testing.T) {
	targetSystem := MultiAgentSystem{}
	manager := manager.NewManager(&targetSystem, 100, false)
	for generationIndex := 0; generationIndex < 100; generationIndex++ {
		manager.SimulateGeneration()
	}

	os.RemoveAll("data")
	os.RemoveAll("logs")
}
