package main

import (
	"math"
	"time"

	"golang.org/x/exp/rand"

	pongsystem "github.com/Otago-Computer-Science-Society/FoosballGeneticLearning/cmd/main/pongSystem"
	geneticbreeder "github.com/Otago-Computer-Science-Society/FoosballGeneticLearning/pkg/GeneticBreeder"
	manager "github.com/Otago-Computer-Science-Society/FoosballGeneticLearning/pkg/Manager"
)

func main() {
	targetSystem := pongsystem.NewPongSystem()

	geneticBreeder := geneticbreeder.NewGeneticBreeder(
		rand.NewSource(uint64(time.Now().Nanosecond())),
		[]float64{0.0, 0.0, 1.0, 1.0},
		[]float64{0.0, 1.0},
		1,
		math.Pow10(-6))
	manager.SimulateManyGenerations(50, 10)
	manager := manager.NewManager(targetSystem, 100, 16, geneticBreeder, true)
	manager.WriteStop()
}
