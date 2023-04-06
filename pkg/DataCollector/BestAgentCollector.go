package datacollector

import (
	"github.com/xitongsys/parquet-go/source"
	"github.com/xitongsys/parquet-go/writer"
)

type BestAgentData struct {
	Score               float64 `parquet:"name=Score, type=DOUBLE"`
	FlatAgentChromosome float64 `parquet:"name=FlatAgentChromosome, type=DOUBLE, repetitiontype=REPEATED"`
}

type BestAgentDataCollector struct {
	dataWriter *writer.ParquetWriter
	fileHandle *source.ParquetFile
}

func NewBestAgentDataCollector(dataFile string) *BestAgentDataCollector {
	fileHandle, dataWriter := newParquetWriter(dataFile, new(GenerationEndData))
	return &BestAgentDataCollector{
		dataWriter: dataWriter,
		fileHandle: fileHandle,
	}
}

func (dc *BestAgentDataCollector) CollectBestAgentData(data GenerationEndData) {
	dc.dataWriter.Write(data)
}
