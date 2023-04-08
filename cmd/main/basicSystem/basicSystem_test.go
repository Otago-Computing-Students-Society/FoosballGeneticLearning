package basicsystem

import (
	"os"
	"testing"

	manager "github.com/Otago-Computer-Science-Society/Foosball-Genetic-Learning/pkg/Manager"
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
