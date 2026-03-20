package dom

import "iter"

// NodeList is an ordered, unique collection of Nodes with fast ID lookup.
type NodeList struct {
	nodes   []Node
	byID    map[string]Node
	set     map[Node]bool
	updates int
}

// NewNodeList returns an empty NodeList.
func NewNodeList() *NodeList {
	return &NodeList{
		byID: make(map[string]Node),
		set:  make(map[Node]bool),
	}
}

func (nl *NodeList) Len() int              { return len(nl.nodes) }
func (nl *NodeList) Get(i int) Node        { return nl.nodes[i] }
func (nl *NodeList) Contains(node Node) bool { return nl.set[node] }
func (nl *NodeList) GetByID(id string) Node  { return nl.byID[id] }
func (nl *NodeList) Updates() int          { return nl.updates }

func (nl *NodeList) Append(node Node) {
	if nl.set[node] {
		return
	}
	id := node.NodeID()
	if id != "" {
		if existing, ok := nl.byID[id]; ok && existing != node {
			nl.remove(existing)
		}
		nl.byID[id] = node
	}
	nl.nodes = append(nl.nodes, node)
	nl.set[node] = true
	nl.updates++
}

func (nl *NodeList) Insert(i int, node Node) {
	if i < 0 {
		i = 0
	}
	if i >= len(nl.nodes) {
		nl.Append(node)
		return
	}
	if nl.set[node] {
		return
	}
	id := node.NodeID()
	if id != "" {
		if existing, ok := nl.byID[id]; ok && existing != node {
			nl.remove(existing)
		}
		nl.byID[id] = node
	}
	nl.nodes = append(nl.nodes, nil)
	copy(nl.nodes[i+1:], nl.nodes[i:])
	nl.nodes[i] = node
	nl.set[node] = true
	nl.updates++
}

func (nl *NodeList) Remove(node Node) bool {
	if !nl.set[node] {
		return false
	}
	nl.remove(node)
	nl.updates++
	return true
}

func (nl *NodeList) remove(node Node) {
	for i, n := range nl.nodes {
		if n == node {
			nl.nodes = append(nl.nodes[:i], nl.nodes[i+1:]...)
			break
		}
	}
	delete(nl.set, node)
	id := node.NodeID()
	if id != "" {
		if nl.byID[id] == node {
			delete(nl.byID, id)
		}
	}
}

func (nl *NodeList) Clear() {
	nl.nodes = nl.nodes[:0]
	nl.byID = make(map[string]Node)
	nl.set = make(map[Node]bool)
	nl.updates++
}

func (nl *NodeList) Index(node Node) int {
	for i, n := range nl.nodes {
		if n == node {
			return i
		}
	}
	return -1
}

// Displayed returns all nodes for which Display() returns true.
func (nl *NodeList) Displayed() []Node {
	result := make([]Node, 0, len(nl.nodes))
	for _, n := range nl.nodes {
		if n.Display() {
			result = append(result, n)
		}
	}
	return result
}

// Iter returns an iterator over all nodes in order.
func (nl *NodeList) Iter() iter.Seq[Node] {
	return func(yield func(Node) bool) {
		for _, n := range nl.nodes {
			if !yield(n) {
				return
			}
		}
	}
}

// Slice returns a copy of the underlying node slice.
func (nl *NodeList) Slice() []Node {
	result := make([]Node, len(nl.nodes))
	copy(result, nl.nodes)
	return result
}
