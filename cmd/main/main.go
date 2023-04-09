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
		[]float64{0.0, 0.0, 1.0, 1.0, 1.0},
		5,
		math.Pow10(-6))
	manager := manager.NewManager(targetSystem, 1000, geneticBreeder, true)
	manager.SimulateManyGenerations(100)
	manager.WriteStop()
}
