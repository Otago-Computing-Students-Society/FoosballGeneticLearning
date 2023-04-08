package datacollector

import (
	"path"

	agent "github.com/Otago-Computer-Science-Society/Foosball-Genetic-Learning/pkg/Agent"

	"github.com/hmcalister/gonum-matrix-io/pkg/gonumio"
	"github.com/xitongsys/parquet-go/source"
	"github.com/xitongsys/parquet-go/writer"
)

const (
	bestAgentDataFile       = "bestAgentData.pq"
	bestAgentChromosomeFile = "bestAgentChromosome.bin"
)

type bestAgentData struct {
	Score float64 `parquet:"name=Score, type=DOUBLE"`
}

type BestAgentDataCollector struct {
	dataDirectory string
	dataWriter    *writer.ParquetWriter
	fileHandle    *source.ParquetFile
}

// Create a new BestAgentDataCollector for storing information on the best agent in a generation
func NewBestAgentDataCollector(dataDirectory string) *BestAgentDataCollector {
	fileHandle, dataWriter := newParquetWriter(path.Join(dataDirectory, bestAgentDataFile), new(bestAgentData))
	return &BestAgentDataCollector{
		dataDirectory: dataDirectory,
		dataWriter:    dataWriter,
		fileHandle:    fileHandle,
	}
}

// Save all relevant information about the best agent to a parquet file
// as well as saving the best agent chromosome to a binary file
func (dc *BestAgentDataCollector) CollectBestAgentData(bestAgent *agent.Agent) {
	dc.dataWriter.Write(bestAgentData{
		Score: bestAgent.Score,
	})
	gonumio.SaveMatrix(bestAgent.Chromosome, path.Join(dc.dataDirectory, bestAgentChromosomeFile))
}
