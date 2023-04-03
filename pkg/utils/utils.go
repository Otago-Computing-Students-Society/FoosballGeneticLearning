package utils

func IsElementInSlice[T comparable](slice []T, element T) bool {
	for _, sliceElement := range slice {
		if sliceElement == element {
			return true
		}
	}
	return false
}
