package manager

import (
	agent "OCSS/FoosballGeneticLearning/pkg/Agent"
	datacollector "OCSS/FoosballGeneticLearning/pkg/DataCollector"
	geneticbreeder "OCSS/FoosballGeneticLearning/pkg/GeneticBreeder"
	simulator "OCSS/FoosballGeneticLearning/pkg/Simulator"
	system "OCSS/FoosballGeneticLearning/pkg/System"
	"io"
	"log"
	"os"
	"path"
	"sync"
	"time"

	"golang.org/x/exp/rand"
)

const (
	DATA_DIRECTORY = "data/"
	LOG_FILE_PATH  = "logs/log"
)

type Manager struct {
	system                      system.System
	logger                      *log.Logger
	generationIndex             int
	numSimulationsPerGeneration int
	currentGeneration           []*agent.Agent
	geneticBreeder              *geneticbreeder.GeneticBreeder
	bestAgentDataCollector      *datacollector.BestAgentDataCollector
	generationEndDataCollector  *datacollector.GenerationEndDataCollector
}

// Create a new manager given the system that is to be learned, and the number of simulations to run per generation
//
// Note the system given must fully implement the System interface in `pkg/system`
func NewManager(system system.System, numSimulationsPerGeneration int) *Manager {
	os.MkdirAll(path.Dir(DATA_DIRECTORY), 0700)
	os.MkdirAll(path.Dir(LOG_FILE_PATH), 0700)

	logFile, err := os.Create(LOG_FILE_PATH)
	if err != nil {
		panic("Could not open log file!")
	}
	logger := log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	return &Manager{
		System:                     system,
		Logger:                     logger,
		bestAgentDataCollector:     datacollector.NewBestAgentDataCollector(DATA_DIRECTORY),
		generationEndDataCollector: datacollector.NewGenerationEndCollector(DATA_DIRECTORY),
	}
}
