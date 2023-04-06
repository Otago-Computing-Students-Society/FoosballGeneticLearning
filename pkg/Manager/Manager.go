package manager

import (
	datacollector "OCSS/FoosballGeneticLearning/pkg/DataCollector"
	system "OCSS/FoosballGeneticLearning/pkg/System"
	"log"
	"os"
	"path"
)

const (
	DATA_DIRECTORY = "data/"
	LOG_FILE_PATH  = "logs/log"
)

type Manager struct {
	System                     system.System
	Logger                     *log.Logger
	bestAgentDataCollector     *datacollector.BestAgentDataCollector
	generationEndDataCollector *datacollector.GenerationEndDataCollector
}

// Create a new manager given the system that is to be learned
//
// Note the system given must fully implement the System interface in `pkg/system`
func NewManager(system system.System) *Manager {
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
