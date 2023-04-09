package utils

import (
	"math"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/rand"
)

func IsElementInSlice[T comparable](slice []T, element T) bool {
	for _, sliceElement := range slice {
		if sliceElement == element {
			return true
		}
	}
	return false
}

// Returns the minimum element of a slice
//
// Panics if slice is empty (len = 0)
func MinElementInSlice[T constraints.Ordered](slice []T) T {
	if len(slice) == 0 {
		panic("cannot find minimum of empty slice")
	}

	currentMin := slice[0]
	for _, elem := range slice {
		if elem < currentMin {
			currentMin = elem
		}
	}

	return currentMin
}

// Returns the maximum element of a slice
//
// Panics if slice is empty (len = 0)
func MaxElementInSlice[T constraints.Ordered](slice []T) T {
	if len(slice) == 0 {
		panic("cannot find minimum of empty slice")
	}

	currentMax := slice[0]
	for _, elem := range slice {
		if elem > currentMax {
			currentMax = elem
		}
	}

	return currentMax
}

// Shuffles the given list
func ShuffleSlice[T comparable](randomGenerator *rand.Rand, slice []T) {
	randomGenerator.Shuffle(len(slice), func(i int, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
}

// Finds the mean and standard deviation of a slice of floats
//
// # Returns
//
// (mean, std)
func SummaryStatistics(slice []float64) (float64, float64) {
	sum := 0.0
	for _, i := range slice {
		sum += i
	}
	mean := sum / float64(len(slice))

	squareDifferenceSum := 0.0
	for _, i := range slice {
		squareDifferenceSum += math.Pow(i-mean, 2.0)
	}
	std := math.Sqrt(squareDifferenceSum / float64(len(slice)))

	return mean, std
}

// Create a new parquet writer to a given file path, using a given struct.
//
// This is a utility method to avoid the same boilerplate code over and over.
//
// See this example (https://github.com/xitongsys/parquet-go/blob/master/example/local_flat.go)
// for information on how to format the structs and use this method nicely.
//
// It may be wise to call `defer writer.WriteStop()` after calling this method!
//
// # Arguments
//
// dataFilePath string: The path to the data file required
//
// dataStruct (generic struct): A valid struct for writing in the parquet format.
// Should be called with `new(dataStruct)` as argument.
//
// # Returns
//
// A ParquetWriter to the data file in question.
func NewParquetWriter[T interface{}](dataFilePath string, dataStruct T) (*source.ParquetFile, *writer.ParquetWriter) {
	os.Remove(dataFilePath)
	dataFileWriter, _ := local.NewLocalFileWriter(dataFilePath)
	parquetDataWriter, _ := writer.NewParquetWriter(dataFileWriter, dataStruct, 4)
	parquetDataWriter.RowGroupSize = 128 * 1024 * 1024 //128MB
	parquetDataWriter.PageSize = 8 * 1024              //8K
	parquetDataWriter.CompressionType = parquet.CompressionCodec_SNAPPY
	parquetDataWriter.Flush(true)

	return &dataFileWriter, parquetDataWriter
}
