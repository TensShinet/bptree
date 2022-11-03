package bptree

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBptree(t *testing.T) {
	tree, _ := New(bptreeMinDegree)

	maxKey := 40

	for i := maxKey; i >= 0; i-- {
		tree.Insert(i, i)
	}

	allKeys := tree.getAllKeys()

	assert.Equal(t, maxKey+1, len(allKeys))

	for i := 0; i <= maxKey; i++ {
		assert.Equal(t, i, allKeys[i])
	}

}
