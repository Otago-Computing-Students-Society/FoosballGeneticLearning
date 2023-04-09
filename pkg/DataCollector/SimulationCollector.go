package datacollector

import (
	"path"

	systemstate "github.com/Otago-Computer-Science-Society/FoosballGeneticLearning/pkg/SystemState"
	"github.com/Otago-Computer-Science-Society/FoosballGeneticLearning/pkg/utils"

	"github.com/xitongsys/parquet-go/source"
	"github.com/xitongsys/parquet-go/writer"
)

type simulationData struct {
	StateIndex  int       `parquet:"name=StateIndex, type=INT32"`
	StateVector []float64 `parquet:"name=StateVector, type=DOUBLE, repetitiontype=REPEATED"`
}

type SimulationDataCollector struct {
	dataDirectory string
	dataFile      string
	dataWriter    *writer.ParquetWriter
	fileHandle    *source.ParquetFile
}

// Create a new BestAgentDataCollector for storing information on the best agent in a generation
func NewSimulationDataCollector(dataDirectory string, dataFile string) *SimulationDataCollector {
	fileHandle, dataWriter := utils.NewParquetWriter(path.Join(dataDirectory, dataFile), new(simulationData))
	return &SimulationDataCollector{
		dataDirectory: dataDirectory,
		dataFile:      dataFile,
		dataWriter:    dataWriter,
		fileHandle:    fileHandle,
	}
}

// Save all relevant information about the best agent to a parquet file
// as well as saving the best agent chromosome to a binary file
func (dc *SimulationDataCollector) CollectSimulationData(state *systemstate.SystemState) {

	dc.dataWriter.Write(simulationData{
		StateIndex:  state.StateIndex,
		StateVector: state.StateVector.RawVector().Data,
	})
}

func (dc *SimulationDataCollector) WriteStop() error {
	if err := dc.dataWriter.WriteStop(); err != nil {
		return err
	}
	if err := (*dc.fileHandle).Close(); err != nil {
		return err
	}
	return nil
}
