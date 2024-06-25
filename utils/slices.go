package utils

func Map[E any, U any](slice []E, fn func(E) U) []U {
	var resultSlice []U
	for _, element := range slice {
		resultSlice = append(resultSlice, fn(element))
	}
	return resultSlice
}
