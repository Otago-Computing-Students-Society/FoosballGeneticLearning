package main

import (
	"math"
	"time"

	"golang.org/x/exp/rand"

	flyingagents "github.com/Otago-Computer-Science-Society/FoosballGeneticLearning/cmd/main/flyingAgents"
	geneticbreeder "github.com/Otago-Computer-Science-Society/FoosballGeneticLearning/pkg/GeneticBreeder"
	manager "github.com/Otago-Computer-Science-Society/FoosballGeneticLearning/pkg/Manager"
)

func main() {
	targetSystem := flyingagents.NewFlyingAgentSystem()

	geneticBreeder := geneticbreeder.NewGeneticBreeder(
		rand.NewSource(uint64(time.Now().Nanosecond())),
		[]float64{0.0, 0.0, 1.0, 1.0, 1.0},
		[]float64{0.0, 1.0, 1.0, 1.0},
		1,
		math.Pow10(-6))
	manager := manager.NewManager(targetSystem, 2500, 10, 16, geneticBreeder, true)
	manager.SimulateManyGenerations(50)
	manager.WriteStop()
}
