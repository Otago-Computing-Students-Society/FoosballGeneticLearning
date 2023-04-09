package main

import (
	"math"
	"time"

	"golang.org/x/exp/rand"

	sys "github.com/Otago-Computer-Science-Society/FoosballGeneticLearning/cmd/main/multiAgentSystem"
	geneticbreeder "github.com/Otago-Computer-Science-Society/FoosballGeneticLearning/pkg/GeneticBreeder"
	manager "github.com/Otago-Computer-Science-Society/FoosballGeneticLearning/pkg/Manager"
)

func main() {
	targetSystem := sys.MultiAgentSystem{}

	geneticBreeder := geneticbreeder.NewGeneticBreeder(
		rand.NewSource(uint64(time.Now().Nanosecond())),
		[]float64{0.0, 0.0, 0.5, 0.5},
		[]float64{0.0, 0.0, 0.0, 0.2, 0.2, 0.2, 0.2, 0.2},
		math.Pow10(-6))
	manager := manager.NewManager(&targetSystem, 100, geneticBreeder, true)
	for generationIndex := 0; generationIndex < 1000; generationIndex++ {
		manager.SimulateGeneration()
	}
}
