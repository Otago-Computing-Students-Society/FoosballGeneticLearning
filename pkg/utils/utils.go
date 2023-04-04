package utils

import (
	"math"

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

// Shuffles the given list
func ShuffleSlice[T comparable](randomGenerator *rand.Rand, slice []T) {
	randomGenerator.Shuffle(len(slice), func(i int, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
}

