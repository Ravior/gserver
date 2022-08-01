package gtimer

import (
	"fmt"
	"testing"
)

func Test_PriorityQueue(t *testing.T) {
	queue := newPriorityQueue()
	queue.Push(1, 20)
	queue.Push(2, 10)
	queue.Push(3, 30)
	fmt.Println(queue.Pop())
	fmt.Println(queue.Pop())
	fmt.Println(queue.Pop())
}
