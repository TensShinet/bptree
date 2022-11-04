package bptree

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
