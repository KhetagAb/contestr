package util

func Transform[T1 any, T2 any](elements []T1, transform func(T1) T2) []T2 {
	result := make([]T2, 0, len(elements))
	for _, element := range elements {
		result = append(result, transform(element))
	}
	return result
}

func Flatten[T any](slices [][]T) []T {
	result := make([]T, 0, len(slices))
	for _, slice := range slices {
		result = append(result, slice...)
	}
	return result
}
