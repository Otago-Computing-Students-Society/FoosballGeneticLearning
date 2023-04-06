package datacollector

import (
	"path"

	"github.com/xitongsys/parquet-go/source"
	"github.com/xitongsys/parquet-go/writer"
)

const (
	bestAgentDataFile = "bestAgentData.pq"
)

type BestAgentData struct {
	Score               float64 `parquet:"name=Score, type=DOUBLE"`
	FlatAgentChromosome float64 `parquet:"name=FlatAgentChromosome, type=DOUBLE, repetitiontype=REPEATED"`
}

type BestAgentDataCollector struct {
	dataWriter *writer.ParquetWriter
	fileHandle *source.ParquetFile
}

func NewBestAgentDataCollector(dataDirectory string) *BestAgentDataCollector {
	fileHandle, dataWriter := newParquetWriter(path.Join(dataDirectory, bestAgentDataFile), new(GenerationEndData))
	return &BestAgentDataCollector{
		dataWriter: dataWriter,
		fileHandle: fileHandle,
	}
}

func (dc *BestAgentDataCollector) CollectBestAgentData(data GenerationEndData) {
	dc.dataWriter.Write(data)
}
