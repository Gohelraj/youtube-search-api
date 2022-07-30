package utils

// GetIndexOf returns the index of the given element in the given slice
func GetIndexOf(element string, data []string) int {
	for index, value := range data {
		if value == element {
			return index
		}
	}
	return -1
}
