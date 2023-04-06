package datacollector

import (
	"github.com/xitongsys/parquet-go/source"
	"github.com/xitongsys/parquet-go/writer"
)

type GenerationEndData struct {
	ScoreMax float64   `parquet:"name=ScoreMax, type=DOUBLE"`
	ScoreMin float64   `parquet:"name=ScoreMin, type=DOUBLE"`
	Scores   []float64 `parquet:"name=Scores, type=DOUBLE, repetitiontype=REPEATED"`
}

type GenerationEndDataCollector struct {
	dataWriter *writer.ParquetWriter
	fileHandle *source.ParquetFile
}

func NewGenerationEndCollector(dataFile string) *GenerationEndDataCollector {
	fileHandle, dataWriter := newParquetWriter(dataFile, new(GenerationEndData))
	return &GenerationEndDataCollector{
		dataWriter: dataWriter,
		fileHandle: fileHandle,
	}
}

func (dc *GenerationEndDataCollector) CollectGenerationEndData(data GenerationEndData) {
	dc.dataWriter.Write(data)
}
