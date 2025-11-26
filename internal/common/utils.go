package common

import (
	"fmt"
)

func GetUniqueValues[S ~[]E, E any, T comparable](sourceSlice S, getValue func(E) T) []T {
	seen := make(map[T]struct{})
	for _, item := range sourceSlice {
		seen[getValue(item)] = struct{}{}
	}

	result := make([]T, 0, len(seen))
	for value := range seen {
		result = append(result, value)
	}
	return result
}

func GetTemplatePath(name string) string {
	return fmt.Sprintf("%s.tmpl", name)
}
