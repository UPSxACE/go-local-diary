package utils

func SafeIndexAccess[T any](slice []T, index int) (value T, error bool) {
	if index >= 0 && index < len(slice) {
		return slice[index], false
	}

	// Return a default value and true to indicate failure
	var zeroVal T
	return zeroVal, true
}