package manager

import (
	"io"
	"log"
	"os"
	"os/signal"
	"path"
	"sync"
	"time"

	agent "github.com/Otago-Computer-Science-Society/FoosballGeneticLearning/pkg/Agent"
	datacollector "github.com/Otago-Computer-Science-Society/FoosballGeneticLearning/pkg/DataCollector"
	geneticbreeder "github.com/Otago-Computer-Science-Society/FoosballGeneticLearning/pkg/GeneticBreeder"
	simulator "github.com/Otago-Computer-Science-Society/FoosballGeneticLearning/pkg/Simulator"
	system "github.com/Otago-Computer-Science-Society/FoosballGeneticLearning/pkg/System"
	"github.com/Otago-Computer-Science-Society/FoosballGeneticLearning/pkg/utils"
	"github.com/schollz/progressbar/v3"
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
	randomGenerator             *rand.Rand
	geneticBreeder              *geneticbreeder.GeneticBreeder
	bestAgentDataCollector      *datacollector.BestAgentDataCollector
	generationEndDataCollector  *datacollector.GenerationEndDataCollector
}

// Create a new manager given the system that is to be learned, and the number of simulations to run per generation
//
// System given must fully implement the System interface in `pkg/system`
// numSimulationsPerGeneration determines how many simulations will be run (and hence the number of agents)
// verbose is a bool flag determining if logs are printed to stdout as well as the log file
func NewManager(system system.System, numSimulationsPerGeneration int, geneticBreeder *geneticbreeder.GeneticBreeder, verbose bool) *Manager {
	os.MkdirAll(path.Dir(DATA_DIRECTORY), 0700)
	os.MkdirAll(path.Dir(LOG_FILE_PATH), 0700)

	logFile, err := os.Create(LOG_FILE_PATH)
	if err != nil {
		panic("Could not open log file!")
	}

	var multiWriter io.Writer
	if verbose {
		multiWriter = io.MultiWriter(os.Stdout, logFile)
	} else {
		multiWriter = io.MultiWriter(logFile)
	}
	logger := log.New(multiWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	numAgents := system.NumAgentsPerSimulation() * numSimulationsPerGeneration
	currentGeneration := make([]*agent.Agent, numAgents)
	for agentIndex := range currentGeneration {
		currentGeneration[agentIndex] = agent.NewRandomGaussianAgent(system.NumActions(), system.NumPercepts())
	}

	randomGenerator := rand.New(rand.NewSource(uint64(time.Now().Nanosecond())))

	return &Manager{
		system:                      system,
		logger:                      logger,
		generationIndex:             0,
		numSimulationsPerGeneration: numSimulationsPerGeneration,
		currentGeneration:           currentGeneration,
		geneticBreeder:              geneticBreeder,
		randomGenerator:             randomGenerator,
		bestAgentDataCollector:      datacollector.NewBestAgentDataCollector(DATA_DIRECTORY),
		generationEndDataCollector:  datacollector.NewGenerationEndCollector(DATA_DIRECTORY),
	}
}

// Simulate a single generation of the system, updating the data writers and breeding the next generation
//
// simulationPerGeneration defines how many simulations to run before tallying up the agents
// scores and breeding a new generation. A large number is better, as it averages agent
// performance.
func (manager *Manager) SimulateGeneration(simulationsPerGeneration int) {
	defer func() { manager.generationIndex += 1 }()
	manager.logger.Printf("STARTING SIMULATION OF GENERATION %v\n", manager.generationIndex)
	numAgentsPerSimulation := manager.system.NumAgentsPerSimulation()

	// Keep a waitgroup to keep track of how many simulations have finished
	var simulationWaitGroup sync.WaitGroup

	// Use a progressbar to track how far through the simulations we are.
	// If we are only doing a single generation, use a silent progress bar instead
	var simulationProgressBar *progressbar.ProgressBar
	if simulationsPerGeneration > 1 {
		simulationProgressBar = progressbar.Default(int64(simulationsPerGeneration), "SIMULATIONS")
	} else {
		simulationProgressBar = progressbar.DefaultSilent(int64(simulationsPerGeneration))
	}

	// Simulate potentially many times
	for simulationIndex := 0; simulationIndex < simulationsPerGeneration; simulationIndex++ {
		// Shuffle the agents (to avoid bias)
		utils.ShuffleSlice(manager.randomGenerator, manager.currentGeneration)
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
		simulationProgressBar.Add(1)
	}
	manager.logger.Println("FINISHED SIMULATING GENERATION")

	// Find the best agent by score
	bestAgent := manager.currentGeneration[0]
	for _, agent := range manager.currentGeneration {
		if bestAgent.Score < agent.Score {
			bestAgent = agent
		}
	}

	// With the best agent, simulate against self (once) and save the result
	manager.logger.Printf("BEST AGENT SCORE: %v\n", bestAgent.Score)
	manager.logger.Printf("SIMULATING BEST AGENT AGAINST SELF")
	simulationDataCollector := datacollector.NewSimulationDataCollector(DATA_DIRECTORY, "BestAgentSimulation.pq")
	simulator.SimulateSystemWithSave(manager.system, []*agent.Agent{bestAgent, bestAgent}, simulationDataCollector)
	simulationDataCollector.WriteStop()
	manager.logger.Printf("FINISHED BEST AGENT SIMULATION")

	// Put data into parquet files
	manager.generationEndDataCollector.CollectGenerationEndData(manager.currentGeneration)
	manager.bestAgentDataCollector.CollectBestAgentData(bestAgent)

	manager.currentGeneration = manager.geneticBreeder.NextGeneration(manager.currentGeneration)
	manager.logger.Printf("--------------------------------------------------------------------------------")
}

// Simulate many generations at once, with handling for SIGINT
func (manager *Manager) SimulateManyGenerations(numGenerations int, simulationsPerGeneration int) {
	sigintChannel := make(chan os.Signal, 1)
	signal.Notify(sigintChannel, os.Interrupt)
simulationGenerationLoop:
	for generationIndex := 0; generationIndex < numGenerations; generationIndex++ {
		select {
		case <-sigintChannel:
			manager.logger.Println("GOT KEYBOARD INTERRUPT")
			break simulationGenerationLoop
		default:
		}
		manager.SimulateGeneration(simulationsPerGeneration)
	}
}

// Flush contents of data collectors to disk and safely close all files
func (manager *Manager) WriteStop() {
	manager.bestAgentDataCollector.WriteStop()
	manager.generationEndDataCollector.WriteStop()
}
