package bptree

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sort"
	"testing"
	"time"
)

func shuffle(slice []int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(slice) > 0 {
		n := len(slice)
		randIndex := r.Intn(n)
		slice[n-1], slice[randIndex] = slice[randIndex], slice[n-1]
		slice = slice[:n-1]
	}
}

func TestBptreeInsert(t *testing.T) {
	tree, _ := New(5)

	maxKey := 10000
	ukSlice := make([]int, 0, maxKey+1)

	for i := 0; i <= maxKey; i++ {
		ukSlice = append(ukSlice, i)
	}

	shuffle(ukSlice)

	for _, v := range ukSlice {
		tree.Insert(v, v)
	}

	allKeys := tree.getAllKeys()

	assert.Equal(t, maxKey+1, len(allKeys))

	for i := 0; i <= maxKey; i++ {
		assert.Equal(t, i, allKeys[i])
	}
}

func remove(slice []int, idx int) []int {
	return append(slice[:idx], slice[idx+1:]...)
}

func TestBptreeDelete(t *testing.T) {
	tree, _ := New(3)

	maxKey := 10000
	targetSlice := make([]int, 0, maxKey+1)

	for i := maxKey; i >= 0; i-- {
		tree.Insert(i, i)
		targetSlice = append(targetSlice, i)
	}

	sort.Slice(targetSlice, func(i, j int) bool {
		return targetSlice[i] < targetSlice[j]
	})

	for i := 0; i < maxKey; i++ {
		dk := rand.Intn(maxKey + 1)
		tree.Delete(dk)
		idx := -1
		for j := 0; j < len(targetSlice); j++ {
			if dk == targetSlice[j] {
				idx = j
				break
			}
		}
		if idx != -1 {
			targetSlice = remove(targetSlice, idx)
		}
	}

	allKeys := tree.getAllKeys()

	assert.Equal(t, len(targetSlice), len(allKeys))

	for i := 0; i < len(targetSlice); i++ {
		assert.Equal(t, targetSlice[i], allKeys[i])
	}
}
