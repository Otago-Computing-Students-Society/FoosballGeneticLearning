package multiagentsystem

import (
	"os"
	"testing"

	manager "github.com/Otago-Computer-Science-Society/Foosball-Genetic-Learning/pkg/Manager"
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
