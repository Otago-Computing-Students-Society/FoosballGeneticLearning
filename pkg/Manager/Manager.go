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

// Simulate a single generation of the system, updating the data writers and breeding the next generation
//
// simulationPerGeneration defines how many simulations to run before tallying up the agents
// scores and breeding a new generation. A large number is better, as it averages agent
// performance.
func (manager *Manager) SimulateGeneration(simulationsPerGeneration int) error {
	defer func() { manager.generationIndex += 1 }()
	sigintChannel := make(chan os.Signal, 1)
	signal.Notify(sigintChannel, os.Interrupt)

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
		select {
		case <-sigintChannel:
			manager.logger.Println("GOT KEYBOARD INTERRUPT")
			return errors.New("got keyboard interrupt")
		default:
		}
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

// Simulate many generations at once, with handling for SIGINT
func (manager *Manager) SimulateManyGenerations(numGenerations int, simulationsPerGeneration int) {
	var err error
	for generationIndex := 0; generationIndex < numGenerations; generationIndex++ {
		err = manager.SimulateGeneration(simulationsPerGeneration)
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
