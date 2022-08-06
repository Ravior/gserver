package gastar

import (
	"testing"
)

func TestList_Add(t *testing.T) {
	nodeList := &list{}
	nodeList.Add(NewNode(1, 1))
	if len(nodeList.nodes) != 1 {
		t.Fail()
	}
}

func TestList_Remove(t *testing.T) {
	nodeList := &list{}
	nodeList.Add(NewNode(1, 1))
	if len(nodeList.nodes) != 1 {
		t.Fail()
	}
	nodeList.Add(NewNode(2, 2))
	nodeList.Add(NewNode(3, 3))
	nodeList.Remove(NewNode(1, 1))
	if len(nodeList.nodes) != 2 {
		t.Fail()
	}
}

func TestList_GetIndex(t *testing.T) {
	nodeList := &list{}
	nodeList.Add(NewNode(1, 1))
	if nodeList.GetIndex(1, 1) != 0 {
		t.Fail()
	}
}

func TestList_ContainsVec(t *testing.T) {
	nodeList := &list{}
	nodeList.Add(NewNode(1, 1))
	nodeList.Add(NewNode(2, 2))
	nodeList.Add(NewNode(3, 3))
	if nodeList.ContainsVec(1, 1) == false {
		t.Fail()
	}
	if nodeList.ContainsVec(3, 4) == true {
		t.Fail()
	}
}

func TestList_Contains(t *testing.T) {
	nodeList := &list{}
	nodeList.Add(NewNode(1, 1))
	nodeList.Add(NewNode(2, 2))
	nodeList.Add(NewNode(3, 3))
	if nodeList.Contains(NewNode(1, 1)) == false {
		t.Fail()
	}
	if nodeList.Contains(NewNode(3, 4)) == true {
		t.Fail()
	}
}

func TestList_IsEmpty(t *testing.T) {
	nodeList := &list{}
	nodeList.Add(NewNode(1, 1))
	nodeList.Add(NewNode(2, 2))
	if nodeList.IsEmpty() == true {
		t.Fail()
	}
	nodeList2 := &list{}
	if nodeList2.IsEmpty() == false {
		t.Fail()
	}
}

func TestList_GetIndexOfMinF(t *testing.T) {
	nodeList := &list{}
	nodeList.Add(&Node{
		g:      0,
		h:      0,
		f:      1,
		X:      0,
		Y:      0,
		parent: nil,
	})
	nodeList.Add(&Node{
		g:      0,
		h:      0,
		f:      2,
		X:      2,
		Y:      2,
		parent: nil,
	})
	if nodeList.GetIndexOfMinF() != 0 {
		t.Fail()
	}

	nodeList.Add(&Node{
		g:      0,
		h:      0,
		f:      0,
		X:      3,
		Y:      3,
		parent: nil,
	})

	if nodeList.GetIndexOfMinF() != 2 {
		t.Fail()
	}
}

func TestList_GetMinFNode(t *testing.T) {
	nodeList := &list{}
	nodeList.Add(&Node{
		g:      0,
		h:      0,
		f:      1,
		X:      1,
		Y:      1,
		parent: nil,
	})
	nodeList.Add(&Node{
		g:      0,
		h:      0,
		f:      2,
		X:      2,
		Y:      2,
		parent: nil,
	})
	if nodeList.GetMinFNode() == nil || nodeList.GetMinFNode().X != 1 {
		t.Fail()
	}

	nodeList.Add(&Node{
		g:      0,
		h:      0,
		f:      0,
		X:      3,
		Y:      3,
		parent: nil,
	})

	if nodeList.GetMinFNode() == nil || nodeList.GetMinFNode().X != 3 {
		t.Fail()
	}
}
