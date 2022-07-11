package v1

func Filter[T any](array []T, filterFunc func(T) bool) []T {
	var filtered []T
	for _, element := range array {
		if filterFunc(element) {
			filtered = append(filtered, element)
		}
	}
	return filtered
}
