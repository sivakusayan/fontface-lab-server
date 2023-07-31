package main

func IndexOf[T comparable](slice []T, value T) int {
	for i, v := range slice {
		if v == value {
			return i
		}
	}
	return -1
}

func Remove[T any](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
}
