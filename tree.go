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
	maxDegree int
	root      *Node
}

func New(maxDegree int) (*Tree, error) {

	if maxDegree < bptreeMinDegree || maxDegree > bptreeMaxDegree {
		return nil, ErrInvalidDegree
	}

	return &Tree{
		maxDegree: maxDegree,
		root:      newLeafNode(maxDegree),
	}, nil
}

func (t *Tree) Find(key int) bool {
	_, idx := t.find(key)
	return idx >= 0
}

func (t *Tree) Insert(key int, value any) {
	n, _ := t.find(key)
	t.insertInLeafNode(n, key, value)

	// Check if this leaf node needs to be split.
	if !n.isFull(t.maxDegree) {
		return
	}
	// Split this leaf node into two leaf nodes, and insert into its parent.
	splitIndex := len(n.keys) / 2
	newNode := newLeafNode(t.maxDegree)
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

func (t *Tree) Delete(key int) {
	n, idx := t.find(key)
	if idx == -1 {
		return
	}

	t.deleteEntry(n, key)
}

func (t *Tree) find(key int) (*Node, int) {
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
			return node, i
		}
	}

	return node, -1
}

func (t *Tree) insertInLeafNode(n *Node, key int, value any) {
	// TODO: binary search?
	i := 0
	for ; i < len(n.keys); i++ {
		if key <= n.keys[i] {
			break
		}
	}

	if i == len(n.keys) {
		n.insertLeafKV(i, key, computeLeafValuesIndex(i), value)
	} else {
		if n.keys[i] == key {
			// Update leaf value.
			n.leafValues[computeLeafValuesIndex(i)] = value
			return
		}
		n.insertLeafKV(i, key, computeLeafValuesIndex(i), value)
	}
}

func (t *Tree) insertInParent(left *Node, key int, right *Node) {
	parent := left.parent
	if parent == nil {
		// The left doesn't have a parent which means The left is the root node before.
		newParent := newNonLeafNode(t.maxDegree)
		insertInSlice(&newParent.keys, len(newParent.keys), key)
		insertInSlice(&newParent.values, len(newParent.values), left)
		insertInSlice(&newParent.values, len(newParent.values), right)

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

	parent.insertNonLeafKV(i, key, computeNonLeafValuesIndex(i), right)

	// Check if this parent node needs to be split.
	if !parent.isFull(t.maxDegree) {
		return
	}

	// Split this leaf node into two leaf nodes, and insert into its parent.
	splitIndex := len(parent.keys) / 2
	splitKey := parent.keys[splitIndex]

	newNode := newNonLeafNode(t.maxDegree)
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

func (t *Tree) deleteEntry(n *Node, key int) {
	keyIdx := getTargetIndex(n.keys, key)
	if n.isLeaf {
		n.deleteLeafKV(keyIdx, computeLeafValuesIndex(keyIdx))
	} else {
		n.deleteNonLeafKV(keyIdx, computeNonLeafValuesIndex(keyIdx))
	}

	if t.root == n && len(n.values) == 1 {
		// This root doesn't contain any keys. Need to choose a new root.
		t.root = n.values[0]
		n.values[0].parent = nil
		return
	}

	if !n.hasTooLessKeys(t.maxDegree) {
		return
	}

	var (
		siblingNode           *Node
		nodeAddrInParentIdx   int
		nodeKeyInParentIdx    int
		isPreviousSiblingNode bool
		middleKey             int

		parent = n.parent
	)

	nodeAddrInParentIdx = getTargetIndex(parent.values, n)
	if nodeAddrInParentIdx == 0 {
		isPreviousSiblingNode = false
		siblingNode = parent.values[nodeAddrInParentIdx+1]
		nodeKeyInParentIdx = nodeAddrInParentIdx
	} else {
		isPreviousSiblingNode = true
		siblingNode = parent.values[nodeAddrInParentIdx-1]
		nodeKeyInParentIdx = nodeAddrInParentIdx - 1
	}

	middleKey = parent.keys[nodeKeyInParentIdx]

	if canFitInOneNode(len(siblingNode.keys), len(n.keys), t.maxDegree, n.isLeaf) {
		// Need to coalesce these two nodes.
		leftNode, rightNode := siblingNode, n
		if !isPreviousSiblingNode {
			leftNode, rightNode = n, siblingNode
		}

		if leftNode.isLeaf {
			leftNode.keys = append(leftNode.keys, rightNode.keys...)
			leftNode.leafValues = append(leftNode.leafValues, rightNode.leafValues...)
			leftNode.nextNode = rightNode.nextNode
		} else {
			// For non-leaf nodes, the first values doesn't have a corresponded key in the node.
			// But the corresponded key of the first values is at the parent node.
			leftNode.keys = append(leftNode.keys, middleKey)
			leftNode.keys = append(leftNode.keys, rightNode.keys...)
			leftNode.values = append(leftNode.values, rightNode.values...)

			for _, child := range leftNode.values {
				child.parent = leftNode
			}
		}
		t.deleteEntry(parent, middleKey)
	} else {
		// Need to redistribution: borrow an entry from its sibling node.
		if isPreviousSiblingNode {
			siblingKeyIdx, siblingValueIdx := siblingNode.lastKeyIdx(), siblingNode.lastValueIdx()
			if n.isLeaf {
				// For the leaf nodes, in this case, borrow (siblingNode.lastKey, siblingNode.lastValue) from sibling node.
				// Due to the change of the first key of the leaf node, it needs to replace the key of the parent node by siblingNode.lastKey.
				n.insertLeafKV(0, siblingNode.keys[siblingKeyIdx], 0, siblingNode.leafValues[siblingValueIdx])
				newK := n.keys[0]
				siblingNode.deleteLeafKV(siblingKeyIdx, siblingValueIdx)
				parent.replaceKey(nodeKeyInParentIdx, newK)
			} else {
				// For the non-leaf nodes, in this case, borrow (middleKey, siblingNode.lastValue) from parent node and sibling node.
				// Due to the change of first key of the non-leaf node, it needs to replace the key of the parent node by siblingNode.lastKey.
				n.insertNonLeafKV(0, middleKey, 0, siblingNode.values[siblingValueIdx])
				siblingNode.values[siblingValueIdx].parent = n
				newK := siblingNode.keys[siblingKeyIdx]
				siblingNode.deleteNonLeafKV(siblingKeyIdx, siblingValueIdx)
				parent.replaceKey(nodeKeyInParentIdx, newK)
			}
		} else {
			siblingKeyIdx, siblingValueIdx := 0, 0
			if n.isLeaf {
				// For the leaf nodes, in this case, borrow (siblingNode.firstKey, siblingNode.firstValue) from the sibling node.
				// Due to the change of the first key of the sibling node, it needs to replace the key of the parent node by siblingNode.firstKey.
				n.insertLeafKV(len(n.keys), siblingNode.keys[siblingKeyIdx], len(n.leafValues), siblingNode.leafValues[siblingValueIdx])
				newK := siblingNode.keys[0]
				siblingNode.deleteLeafKV(siblingKeyIdx, siblingValueIdx)
				parent.replaceKey(nodeKeyInParentIdx, newK)
			} else {
				// For the non-leaf nodes, in this case, borrow (middleKey, siblingNode.firstValue) from the sibling node.
				// Due to the change of the first key of the sibling node, it needs to replace the key of the parent node by siblingNode.firstKey.
				n.insertNonLeafKV(len(n.keys), middleKey, len(n.values), siblingNode.values[siblingValueIdx])
				siblingNode.values[siblingValueIdx].parent = n
				newK := siblingNode.keys[siblingKeyIdx]
				siblingNode.deleteNonLeafKV(siblingKeyIdx, siblingValueIdx)
				parent.replaceKey(nodeKeyInParentIdx, newK)
			}
		}
	}

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
