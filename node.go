package bptree

type Node struct {
	isLeaf bool
	// The max length of keys is maxDegree - 1. While the length of keys
	// is equal to the maxDegree, the node needs to be split into two nodes.
	keys       []int   // To store the search keys.
	parent     *Node   // To store the parent of this node.
	values     []*Node // To store children for non-leaf nodes.
	leafValues []any   // To store values for the leaf nodes.
	nextNode   *Node   // To store the sibling node for the leaf nodes.
}

func newLeafNode(maxDegree int) *Node {
	return &Node{
		isLeaf:     true,
		keys:       make([]int, 0, maxDegree),
		leafValues: make([]interface{}, 0, maxDegree),
	}
}

func newNonLeafNode(maxDegree int) *Node {
	return &Node{
		isLeaf: false,
		keys:   make([]int, 0, maxDegree),
		values: make([]*Node, 0, maxDegree+1),
	}
}

func (n *Node) isFull(maxDegree int) bool {
	return len(n.keys) >= maxDegree
}

func (n *Node) hasTooLessKeys(maxDegree int) bool {
	if n.parent == nil {
		// This current node is the root.
		return false
	}
	return (n.isLeaf && len(n.leafValues) < ceil(maxDegree-1, 2)) ||
		(!n.isLeaf && len(n.values) < ceil(maxDegree, 2))
}

func (n *Node) insertLeafKV(keyIdx int, key int, valueIdx int, value any) {
	insertInSlice(&n.keys, keyIdx, key)
	insertInSlice(&n.leafValues, valueIdx, value)
}

func (n *Node) insertNonLeafKV(keyIdx int, key int, valueIdx int, value *Node) {
	insertInSlice(&n.keys, keyIdx, key)
	insertInSlice(&n.values, valueIdx, value)
}

func (n *Node) deleteLeafKV(keyIdx int, valueIdx int) {
	deleteInSlice(&n.keys, keyIdx)
	deleteInSlice(&n.leafValues, valueIdx)
}

func (n *Node) deleteNonLeafKV(keyIdx int, valueIdx int) {
	deleteInSlice(&n.keys, keyIdx)
	deleteInSlice(&n.values, valueIdx)
}

func (n *Node) lastKeyIdx() int {
	return len(n.keys) - 1
}

func (n *Node) lastValueIdx() int {
	if n.isLeaf {
		return len(n.leafValues) - 1
	}
	return len(n.values) - 1
}

func (n *Node) replaceKey(keyIdx int, key int) {
	n.keys[keyIdx] = key
}
