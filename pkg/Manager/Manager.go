package manager

import (
	"errors"
	"io"
	"log"
	"os"
	"os/signal"
	"path"
	"sort"
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
	numThreads                  int
	geneticBreeder              *geneticbreeder.GeneticBreeder
	bestAgentDataCollector      *datacollector.BestAgentDataCollector
	generationEndDataCollector  *datacollector.GenerationEndDataCollector
}

// Create a new manager given the system that is to be learned, and the number of simulations to run per generation
//
// System given must fully implement the System interface in `pkg/system`
//
// numSimulationsPerGeneration defines how many simulations to run before tallying up the agents
// scores and breeding a new generation. A large number is better, as it averages agent
// performance.
//
// verbose is a bool flag determining if logs are printed to stdout as well as the log file
func NewManager(system system.System, numSimulationsPerGeneration int, numThreads int, geneticBreeder *geneticbreeder.GeneticBreeder, verbose bool) *Manager {
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

	if numThreads <= 0 {
		panic("Number of threads must be a positive integer!")
	}

	randomGenerator := rand.New(rand.NewSource(uint64(time.Now().Nanosecond())))

	return &Manager{
		system:                      system,
		logger:                      logger,
		generationIndex:             0,
		numSimulationsPerGeneration: numSimulationsPerGeneration,
		currentGeneration:           currentGeneration,
		geneticBreeder:              geneticBreeder,
		numThreads:                  numThreads,
		randomGenerator:             randomGenerator,
		bestAgentDataCollector:      datacollector.NewBestAgentDataCollector(DATA_DIRECTORY),
		generationEndDataCollector:  datacollector.NewGenerationEndCollector(DATA_DIRECTORY),
	}
}

// Simulate a single repetition, of which there may be many (always at least one) within a generation
// This method is not exposed publicly. The intention is for users to call SimulateGeneration instead.
func (manager *Manager) simulateRepetition() error {
	numAgentsPerSimulation := manager.system.NumAgentsPerSimulation()

	// Channel to send collections of agents through to simulation goroutines.
	// The number of agents sent at once is equal to the number of agents required
	// for one simulation.
	agentChannel := make(chan []*agent.Agent, manager.numSimulationsPerGeneration)
	// Channel to receive signals (hence generic struct{}) for when a simulation finishes.
	simulationFinishedSignalChannel := make(chan struct{})
	// A simple counter of how many simulations are running
	simulationsRunningCounter := 0

	// Create a number of goroutines to handle as simulations
	for i := 0; i < manager.numThreads; i++ {
		go simulator.ConcurrentSimulationRoutine(manager.system, agentChannel, simulationFinishedSignalChannel)
	}

	// Shuffle the agents (to avoid bias)
	utils.ShuffleSlice(manager.randomGenerator, manager.currentGeneration)

	// Actually start all the simulations
	for simulationIndex := 0; simulationIndex < manager.numSimulationsPerGeneration; simulationIndex++ {
		simulationsRunningCounter += 1
		// Find the agents to be used in this simulation
		simulationAgents := manager.currentGeneration[numAgentsPerSimulation*simulationIndex : numAgentsPerSimulation*(simulationIndex+1)]
		// Send the agents to the simulators - blocks until agents can be taken
		agentChannel <- simulationAgents
	}
	close(agentChannel)

	// Count any finished simulations and decrement the simulation counter
	for range simulationFinishedSignalChannel {
		simulationsRunningCounter -= 1
		if simulationsRunningCounter <= 0 {
			break
		}
	}

	return nil
}

// Simulate a single generation of the system, updating the data writers and breeding the next generation
//
// simulationPerGeneration
func (manager *Manager) SimulateGeneration() error {
	defer func() { manager.generationIndex += 1 }()
	sigintChannel := make(chan os.Signal, 1)
	signal.Notify(sigintChannel, os.Interrupt)

	manager.logger.Printf("STARTING SIMULATION OF GENERATION %v\n", manager.generationIndex)

	// Use a progressbar to track how far through the simulations we are.
	// If we are only doing a single generation, use a silent progress bar instead
	var simulationRepeatsProgressBar *progressbar.ProgressBar
	if manager.numSimulationsPerGeneration > 1 {
		simulationRepeatsProgressBar = progressbar.Default(int64(manager.numSimulationsPerGeneration), "SIMULATION REPETITIONS")
	} else {
		simulationRepeatsProgressBar = progressbar.DefaultSilent(int64(manager.numSimulationsPerGeneration))
	}

	// Simulate as many times as required, passing agents through channel to awaiting goroutines
	for simulationRepeatIndex := 0; simulationRepeatIndex < manager.numSimulationsPerGeneration; simulationRepeatIndex++ {
		// Handle any keyboard interrupts
		select {
		case <-sigintChannel:
			manager.logger.Println("GOT KEYBOARD INTERRUPT")
			return errors.New("got keyboard interrupt")
		default:
		}

		err := manager.simulateRepetition()
		if err != nil {
			return err
		}
		simulationRepeatsProgressBar.Add(1)
	}
	manager.logger.Println("FINISHED SIMULATING GENERATION")

	// Find the best agent by score
	sort.Slice(manager.currentGeneration, func(i, j int) bool {
		return manager.currentGeneration[i].Score < manager.currentGeneration[j].Score
	})
	bestAgent := manager.currentGeneration[len(manager.currentGeneration)-1]

	// With the best agent, simulate and save the result
	manager.logger.Printf("BEST AGENT SCORE: %v\n", bestAgent.Score)
	manager.logger.Printf("SIMULATING BEST AGENT'S ")
	// Get the top n agents, where n is the number of agents needed for the simulation
	bestAgentArray := make([]*agent.Agent, manager.system.NumAgentsPerSimulation())
	for bestAgentIndex := range bestAgentArray {
		bestAgentArray[bestAgentIndex] = manager.currentGeneration[len(manager.currentGeneration)-bestAgentIndex-1]
	}
	// Then simulate these and put data into data collector
	simulationDataCollector := datacollector.NewSimulationDataCollector(DATA_DIRECTORY, "BestAgentSimulation.pq")
	simulator.SimulateSystemWithSave(manager.system, bestAgentArray, simulationDataCollector)
	simulationDataCollector.WriteStop()
	manager.logger.Printf("FINISHED BEST AGENT SIMULATION")

	// Put data into parquet files
	manager.generationEndDataCollector.CollectGenerationEndData(manager.currentGeneration)
	manager.bestAgentDataCollector.CollectBestAgentData(bestAgent)

	manager.currentGeneration = manager.geneticBreeder.NextGeneration(manager.currentGeneration)
	manager.logger.Printf("--------------------------------------------------------------------------------")
	return nil
}

// Simulate many generations in a loop
func (manager *Manager) SimulateManyGenerations(numGenerations int) {
	var err error
	for generationIndex := 0; generationIndex < numGenerations; generationIndex++ {
		err = manager.SimulateGeneration()
		if err != nil {
			break
		}
	}
}

// Flush contents of data collectors to disk and safely close all files
func (manager *Manager) WriteStop() {
	manager.bestAgentDataCollector.WriteStop()
	manager.generationEndDataCollector.WriteStop()
}
