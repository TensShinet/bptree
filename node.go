package bptree

type Node struct {
	maxDegree int // To store the maximum degree of a node.
	isLeaf    bool
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
		maxDegree:  maxDegree,
		keys:       make([]int, 0, maxDegree),
		leafValues: make([]interface{}, 0, maxDegree),
	}
}

func newNonLeafNode(maxDegree int) *Node {
	return &Node{
		isLeaf:    false,
		maxDegree: maxDegree,
		keys:      make([]int, 0, maxDegree),
		values:    make([]*Node, 0, maxDegree+1),
	}
}

func (n *Node) isFull() bool {
	return len(n.keys) >= n.maxDegree
}
