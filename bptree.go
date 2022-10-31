package bptree

type Tree struct {
	MaxDegree int
	Root      *Node
}

type Node struct {
	ValueLength int
	IsLeaf      bool
	Values      []int
	Pointers    []*Node
	Parent      *Node
}

func New(MaxDegree int) *Tree {
	return &Tree{
		MaxDegree: MaxDegree,
		Root: &Node{
			ValueLength: 0,
			IsLeaf:      true,
			Values:      make([]int, MaxDegree-1),
			Pointers:    make([]*Node, MaxDegree),
			Parent:      nil,
		},
	}
}

func (t *Tree) Find(v int) (*Node, bool) {
	node := t.Root

	for !node.IsLeaf {
		// To find the smallest value bigger than the v
		// TODO: binary search
		i := 0
		for ; i < node.ValueLength; i++ {
			if node.Values[i] > v {
				break
			}
		}
		node = node.Pointers[i]
	}

	// Now the node is a Leaf Node. And the v will not in the next sibling node.
	// Because the next sibling first value is bigger than the v.

	// Scan the node values and find the v.
	// TODO: binary search
	for i := 0; i < node.ValueLength; i++ {
		if node.Values[i] == v {
			return node, true
		}
	}

	return node, false
}

func (t *Tree) Insert(v int) {

}

func (t *Tree) Delete(v int) bool {
	return false
}
