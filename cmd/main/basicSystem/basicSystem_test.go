package basicsystem

import (
	"math"
	"os"
	"testing"
	"time"

	"golang.org/x/exp/rand"

	geneticbreeder "github.com/Otago-Computer-Science-Society/FoosballGeneticLearning/pkg/GeneticBreeder"
	manager "github.com/Otago-Computer-Science-Society/FoosballGeneticLearning/pkg/Manager"
)

func TestBasicSystem(t *testing.T) {
	targetSystem := BasicSystem{}
	geneticBreeder := geneticbreeder.NewGeneticBreeder(
		rand.NewSource(uint64(time.Now().Nanosecond())),
		[]float64{0.0, 0.0, 0.5, 0.5},
		[]float64{0.0, 0.0, 0.0, 0.2, 0.2, 0.2, 0.2, 0.2},
		0,
		math.Pow10(-6))
	manager := manager.NewManager(&targetSystem, 100, 8, geneticBreeder, false)
	manager.SimulateManyGenerations(100)

	os.RemoveAll("data")
	os.RemoveAll("logs")
}
