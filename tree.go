package bptree

import (
	"errors"
)

const (
	bptreeMinDegree = 3
	bptreeMaxDegree = 20
)

var ErrInvalidDegree = errors.New("invalid degree")

type Tree struct {
	root *Node
}

func New(maxDegree int) (*Tree, error) {

	if maxDegree < bptreeMinDegree || maxDegree > bptreeMaxDegree {
		return nil, ErrInvalidDegree
	}

	return &Tree{
		root: newLeafNode(maxDegree),
	}, nil
}

func (t *Tree) Find(key int) (*Node, bool) {
	node := t.root

	for !node.isLeaf {
		// Find the smallest key bigger than this key.
		// TODO: binary search?
		i := 0
		for ; i < len(node.keys); i++ {
			if node.keys[i] > key {
				break
			}
		}
		node = node.values[i]
	}

	// Now the node is a Leaf Node. And the key will not in the next sibling node.
	// Because the next sibling first key is bigger than the key.

	// Scan the node keys and find this key.
	// TODO: binary search?
	for i := 0; i < len(node.keys); i++ {
		if node.keys[i] == key {
			return node, true
		}
	}

	return node, false
}

func (t *Tree) Insert(key int, value interface{}) {
	n, _ := t.Find(key)
	t.insertInLeafNode(n, key, value)

	// Check if this leaf node needs to be split.
	if !n.isFull() {
		return
	}
	// Split this leaf node into two leaf nodes, and insert into its parent.
	splitIndex := len(n.keys) / 2
	newNode := newLeafNode(n.maxDegree)
	for _, v := range n.keys[splitIndex:] {
		newNode.keys = append(newNode.keys, v)
	}
	for _, v := range n.leafValues[splitIndex:] {
		newNode.leafValues = append(newNode.leafValues, v)
	}
	newNode.nextNode = n.nextNode
	newNode.parent = n.parent

	n.keys = n.keys[:splitIndex]
	n.leafValues = n.leafValues[:splitIndex]
	n.nextNode = newNode

	t.insertInParent(n, newNode.keys[0], newNode)
}

func (t *Tree) insertInLeafNode(n *Node, key int, value interface{}) {
	// TODO: binary search?
	i := 0
	for ; i < len(n.keys); i++ {
		if key <= n.keys[i] {
			break
		}
	}

	if i == len(n.keys) {
		n.keys = append(n.keys, key)
		n.leafValues = append(n.leafValues, value)
	} else {
		if n.keys[i] == key {
			// Update leaf value.
			n.leafValues[i] = value
			return
		}
		// Insert a key and a value to their slices separately to make two empty slots.
		n.keys = append(n.keys, 0)
		n.leafValues = append(n.leafValues, value)

		for j := len(n.leafValues) - 2; j >= i; j-- {
			n.keys[j+1], n.leafValues[j+1] = n.keys[j], n.leafValues[j]
		}

		// Insert the new key and the new value.
		n.keys[i], n.leafValues[i] = key, value
	}
}

func (t *Tree) insertInParent(left *Node, key int, right *Node) {
	parent := left.parent
	if parent == nil {
		// The left doesn't have a parent which means The left is the root node before.
		newParent := newNonLeafNode(left.maxDegree)
		newParent.keys = append(newParent.keys, key)
		newParent.values = append(newParent.values, left, right)

		// Update the left and the right's parents
		left.parent, right.parent = newParent, newParent
		t.root = newParent
		return
	}

	// Insert in parent first, and to check if the parent needs to be split into two nodes
	i := 0
	for ; i < len(parent.keys); i++ {
		if key <= parent.keys[i] {
			break
		}
	}

	if i == len(parent.keys) {
		parent.keys = append(parent.keys, key)
		parent.values = append(parent.values, right)
	} else {
		// Insert a key and a value to their slices separately to make two empty slots.
		parent.keys = append(parent.keys, 0)
		parent.values = append(parent.values, nil)

		for j := len(parent.keys) - 2; j >= i; j-- {
			// Notice: the length of values slice of non-leaf node is always one more than the length of keys slice
			parent.keys[j+1], parent.values[j+2] = parent.keys[j], parent.values[j+1]
		}

		// Insert the new key and the new value.
		parent.keys[i], parent.values[i+1] = key, right
	}

	// Check if this parent node needs to be split.
	if !parent.isFull() {
		return
	}

	// Split this leaf node into two leaf nodes, and insert into its parent.
	splitIndex := len(parent.keys) / 2
	splitKey := parent.keys[splitIndex]

	newNode := newNonLeafNode(parent.maxDegree)
	for _, v := range parent.keys[splitIndex+1:] {
		newNode.keys = append(newNode.keys, v)
	}
	for _, v := range parent.values[splitIndex+1:] {
		newNode.values = append(newNode.values, v)
	}
	newNode.parent = parent.parent

	parent.keys = parent.keys[:splitIndex]
	parent.values = parent.values[:splitIndex+1]

	// Update all children's parents because the parents of some children are wrong.
	for _, child := range parent.values {
		child.parent = parent
	}

	for _, child := range newNode.values {
		child.parent = newNode
	}

	t.insertInParent(parent, splitKey, newNode)
}

func (t *Tree) getAllKeys() []int {
	node := t.root
	for !node.isLeaf {
		node = node.values[0]
	}
	res := make([]int, 0)
	for node.nextNode != nil {
		res = append(res, node.keys...)
		node = node.nextNode
	}

	res = append(res, node.keys...)
	return res
}