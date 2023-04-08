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

// Simulate a single generation of the system, updating the data writers and breeding the next generation
func (manager *Manager) SimulateGeneration() {
	defer func() { manager.generationIndex += 1 }()
	manager.logger.Printf("STARTING SIMULATION OF GENERATION %v\n", manager.generationIndex)
	numAgentsPerSimulation := manager.system.NumAgentsPerSimulation()

	// Keep a waitgroup to keep track of how many simulations have finished
	var simulationWaitGroup sync.WaitGroup
	// Actually start all the simulations
	for simulationIndex := 0; simulationIndex < manager.numSimulationsPerGeneration; simulationIndex++ {
		simulationWaitGroup.Add(1)
		// Find the agents to be used in this simulation
		simulationAgents := manager.currentGeneration[numAgentsPerSimulation*simulationIndex : numAgentsPerSimulation*(simulationIndex+1)]
		// Start simulation in very simple anonymous wrapper - to decrement waitgroup when simulation is done
		go func() {
			defer simulationWaitGroup.Done()
			simulator.SimulateSystem(manager.system, simulationAgents)
		}()
	}
	simulationWaitGroup.Wait()
	manager.logger.Println("FINISHED SIMULATING GENERATION")

	// Find the best agent by score
	bestAgent := manager.currentGeneration[0]
	for _, agent := range manager.currentGeneration {
		if bestAgent.Score < agent.Score {
			bestAgent = agent
		}
	}
	manager.logger.Printf("BEST AGENT SCORE: %v\n", bestAgent.Score)

	// Put these scores into parquet files
	manager.generationEndDataCollector.CollectGenerationEndData(manager.currentGeneration)
	manager.bestAgentDataCollector.CollectBestAgentData(bestAgent)

	manager.currentGeneration = manager.geneticBreeder.NextGeneration(manager.currentGeneration)
	manager.logger.Printf("--------------------------------------------------------------------------------")
}
