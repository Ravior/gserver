package gastar

// list represents a list of nodes
type list struct {
	nodes []*Node
}

// Add one or more nodes to the list
func (l *list) Add(nodes ...*Node) {
	l.nodes = append(l.nodes, nodes...)
}

// Remove a node from the list
// if the node is not found we do nothing
func (l *list) Remove(removeNode *Node) {
	index := l.GetIndex(removeNode.X, removeNode.Y)
	if index >= 0 {
		l.nodes = append(l.nodes[:index], l.nodes[index+1:]...)
	}
}

// GetIndex returns the index of the node in the list
// if the node is not found the return value is -1
func (l *list) GetIndex(x, y int32) int {
	for index, node := range l.nodes {
		if node.X == x && node.Y == y {
			return index
		}
	}
	return -1
}

// ContainsVec check if a vec is in the list
func (l *list) ContainsVec(x, y int32) bool {
	return l.GetIndex(x, y) >= 0
}

// Contains check if a node is in the list
func (l *list) Contains(searchNode *Node) bool {
	return l.GetIndex(searchNode.X, searchNode.Y) >= 0
}

// IsEmpty returns if the nodes list has nodes or not
func (l *list) IsEmpty() bool {
	return len(l.nodes) == 0
}

// Clear removes all nodes from the list
func (l *list) Clear() {
	l.nodes = []*Node{}
}

// GetIndexOfMinF returns the index of the nodes list
// with the smallest node.F value
//
// if no node is found it returns -1
func (l *list) GetIndexOfMinF() int {
	lastNodeIndex := -1
	if len(l.nodes) > 0 {
		lastNode := &Node{}
		for index, node := range l.nodes {
			if lastNodeIndex == -1 || node.f < lastNode.f {
				lastNode = node
				lastNodeIndex = index
			}
		}
	}
	return lastNodeIndex
}

// GetMinFNode returns the node with the smallest node.F value
func (l *list) GetMinFNode() *Node {
	minFIndex := l.GetIndexOfMinF()
	if minFIndex == -1 {
		return nil
	}
	return l.nodes[minFIndex]
}
