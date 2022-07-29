package gastar

type Node struct {
	// 到起点的步数
	g int32
	// 到终点的步数
	h int32
	// f = g+h
	f int32
	// X坐标
	X int32
	// Y坐标
	Y int32
	// 父节点
	parent *Node
}

// NewNode 构造函数
func NewNode(x, y int32) *Node {
	return &Node{
		g:      0,
		h:      0,
		f:      0,
		X:      x,
		Y:      y,
		parent: nil,
	}
}

// 判断两点是否相等
func (n *Node) equal(other *Node) bool {
	return n.X == other.X && n.Y == other.Y
}
