package gastar

// 估算函数：F(x) = G(x) + H(x)
// F：可以认为是起点到终点的全部代价（长度）
// G：从起点到当前点的代价（长度）
// H：从当前点到终点的估算代价（长度），使用的是曼哈顿距离算法来估算的
// OPEN列表：参与估算的点组成的集合（路径中所有 有可能经过的点）
// CLOSE列表：不需要参与估算的点组成的集合
// 曼哈顿距离算法：用来估算两点距离，使用公式length = abs(A.x - B.x) + abs(A.y - B.y)

// AStar 采用左上角坐标系
type AStar struct {
	width       int32                 // 地图宽
	height      int32                 // 地图高
	closedList  list                  // closedList
	openList    list                  // openList
	startNode   *Node                 // 开始寻路节点
	endNode     *Node                 // 结束寻路节点
	canPassFunc func(x, y int32) bool // Node 能否走的判断函数
}

// NewAStar 构造函数
func NewAStar(width, height int32, canPassFunc func(x, y int32) bool) *AStar {
	return &AStar{
		width:       width,
		height:      height,
		canPassFunc: canPassFunc,
	}
}

// 计算绝对值
func (a *AStar) abs(x int32) int32 {
	if x > 0 {
		return x
	} else {
		return -x
	}
}

// h 计算两点之间的距离(曼哈顿距离算法)
// nodeA and nodeB calculates by the manhattan distance
func (a *AStar) h(nodeA *Node, nodeB *Node) int32 {
	return a.abs(nodeA.X-nodeB.X) + a.abs(nodeA.Y-nodeB.Y)
}

// getNeighborNodes calculates the next neighbors of the given node
// if a neighbor node is not accessible the node will be ignored
func (a *AStar) getNeighborNodes(node *Node) []*Node {
	var neighborNodes []*Node

	// 右
	if a.isAccessible(node.X+1, node.Y) {
		neighborNodes = append(neighborNodes, &Node{X: node.X + 1, Y: node.Y, parent: node})
	}

	// 左
	if a.isAccessible(node.X-1, node.Y) {
		neighborNodes = append(neighborNodes, &Node{X: node.X - 1, Y: node.Y, parent: node})
	}

	// 下
	if a.isAccessible(node.X, node.Y+1) {
		neighborNodes = append(neighborNodes, &Node{X: node.X, Y: node.Y + 1, parent: node})
	}

	// 上
	if a.isAccessible(node.X, node.Y-1) {
		neighborNodes = append(neighborNodes, &Node{X: node.X, Y: node.Y - 1, parent: node})
	}

	return neighborNodes
}

// isAccessible checks if the node is reachable in the grid
// and is not in the invalidNodes slice
func (a *AStar) isAccessible(x, y int32) bool {
	// 如果是终点，则直接判定为可以达到目标点
	if a.endNode.X == x && a.endNode.Y == y {
		return true
	}

	// if node is out of bound
	if x >= 0 && x < a.width && y >= 0 && y < a.height {
		if a.closedList.ContainsVec(x, y) {
			return false
		}

		return a.canPassFunc(x, y)
	}
	return false
}

// 判断是否是结束节点
func (a *AStar) isEndNode(checkNode, endNode *Node) bool {
	return checkNode.X == endNode.X && checkNode.Y == endNode.Y
}

// Clear 清理
func (a *AStar) Clear() {
	a.openList.Clear()
	a.closedList.Clear()
}

// Find 寻路
func (a *AStar) Find(startNode, endNode *Node) []*Node {
	// 如果两点相同，则不进行巡路
	if startNode.equal(endNode) {
		return nil
	}

	a.startNode = startNode
	a.endNode = endNode

	defer func() {
		a.Clear()
	}()

	a.openList.Add(startNode)

	for !a.openList.IsEmpty() {

		currentNode := a.openList.GetMinFNode()
		if currentNode == nil {
			return nil
		}

		a.openList.Remove(currentNode)
		a.closedList.Add(currentNode)

		// 找到终点
		if a.isEndNode(currentNode, endNode) {
			return a.getNodePath(currentNode)
		}

		neighbors := a.getNeighborNodes(currentNode)

		for _, neighbor := range neighbors {
			if a.closedList.Contains(neighbor) {
				continue
			}

			a.calculateNode(neighbor)

			if !a.openList.Contains(neighbor) {
				a.openList.Add(neighbor)
			}
		}

	}

	return nil
}

// calculateNode calculates the F, G and h value for the given node
func (a *AStar) calculateNode(node *Node) {
	node.g++
	node.h = a.h(node, a.endNode)
	node.f = node.g + node.h
}

// getNodePath returns the chain of parent nodes
// the given node will be still included in the nodes slice
func (a *AStar) getNodePath(currentNode *Node) []*Node {
	var nodePath []*Node
	nodePath = append(nodePath, currentNode)
	for {
		if currentNode.parent == nil {
			break
		}

		parentNode := currentNode.parent

		// if the end of node chain
		if parentNode.parent == nil {
			break
		}

		nodePath = append(nodePath, parentNode)
		currentNode = parentNode
	}

	// 反转（实际查出来的结果倒序)
	if len(nodePath) > 0 {
		for i, j := 0, len(nodePath)-1; i < j; i, j = i+1, j-1 {
			nodePath[i], nodePath[j] = nodePath[j], nodePath[i]
		}
	}

	return nodePath
}
