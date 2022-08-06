package gastar

import "testing"

func TestNewNode(t *testing.T) {
	node1 := NewNode(1, 1)
	node2 := NewNode(1, 1)
	if node1 == nil || node2 == nil || node1.equal(node2) == false {
		t.Fail()
	}
}
