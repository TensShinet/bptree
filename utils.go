package bptree

import (
	"math"
)

func insertInSlice[T any](s *[]T, idx int, value T) {
	var dummy T
	*s = append(*s, dummy)

	for i := len(*s) - 2; i >= idx; i-- {
		(*s)[i+1] = (*s)[i]
	}

	(*s)[idx] = value
}

func deleteInSlice[T any](s *[]T, idx int) {
	for i := idx + 1; i < len(*s); i++ {
		(*s)[i-1] = (*s)[i]
	}

	*s = (*s)[:len(*s)-1]
}

func getTargetIndex[T comparable](s []T, t T) int {
	for i, v := range s {
		if v == t {
			return i
		}
	}
	return -1
}

func computeLeafValuesIndex(keyIdx int) int {
	return keyIdx
}

func computeNonLeafValuesIndex(keyIdx int) int {
	return keyIdx + 1
}

func ceil(numerator, denominator int) int {
	return int(math.Ceil(float64(numerator) / float64(denominator)))
}

func canFitInOneNode(l1, l2, degree int, isLeaf bool) bool {
	if isLeaf {
		return l1+l2 < degree
	}
	return l1+l2 < degree-1
}
