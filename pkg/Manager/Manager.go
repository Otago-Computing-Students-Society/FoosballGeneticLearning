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
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	logger := log.New(multiWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	numAgents := system.NumAgentsPerSimulation() * numSimulationsPerGeneration
	currentGeneration := make([]*agent.Agent, numAgents)
	for agentIndex := range currentGeneration {
		currentGeneration[agentIndex] = agent.NewRandomGaussianAgent(system.NumActions(), system.NumPercepts())
	}

	geneticBreeder := geneticbreeder.NewGeneticBreeder(rand.NewSource(uint64(time.Now().Nanosecond())))

	return &Manager{
		system:                      system,
		logger:                      logger,
		generationIndex:             0,
		numSimulationsPerGeneration: numSimulationsPerGeneration,
		currentGeneration:           currentGeneration,
		geneticBreeder:              geneticBreeder,
		bestAgentDataCollector:      datacollector.NewBestAgentDataCollector(DATA_DIRECTORY),
		generationEndDataCollector:  datacollector.NewGenerationEndCollector(DATA_DIRECTORY),
	}
}
