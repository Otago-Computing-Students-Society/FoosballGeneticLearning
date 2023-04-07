package datacollector

import (
	agent "OCSS/FoosballGeneticLearning/pkg/Agent"
	"path"

	"github.com/xitongsys/parquet-go/source"
	"github.com/xitongsys/parquet-go/writer"
)

const (
	generationEndDataFile = "generationEndData.pq"
)

type generationEndData struct {
	Scores []float64 `parquet:"name=Scores, type=DOUBLE, repetitiontype=REPEATED"`
}

type GenerationEndDataCollector struct {
	dataWriter *writer.ParquetWriter
	fileHandle *source.ParquetFile
}

func NewGenerationEndCollector(dataDirectory string) *GenerationEndDataCollector {
	fileHandle, dataWriter := newParquetWriter(path.Join(dataDirectory, generationEndDataFile), new(generationEndData))
	return &GenerationEndDataCollector{
		dataWriter: dataWriter,
		fileHandle: fileHandle,
	}
}

func (dc *GenerationEndDataCollector) CollectGenerationEndData(agents []*agent.Agent) {
	scores := make([]float64, len(agents))
	for agentIndex := range agents {
		scores[agentIndex] = agents[agentIndex].Score
	}

	dc.dataWriter.Write(generationEndData{
		Scores: scores,
	})
}
