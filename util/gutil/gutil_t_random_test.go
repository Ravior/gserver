package gutil

import (
	"fmt"
	"testing"
)

func Test_Random(t *testing.T) {
	m := make(map[float64]int)
	for i := 0; i < 10000; i++ {
		// r > 0 && r < 1
		r := Random()
		m[r] = m[r] + 1
	}
	for _, v := range m {
		if v > 1 {
			t.Fail()
		}
	}
}

func Test_RandomSeq(t *testing.T) {
	m := make(map[string]int)
	for i := 0; i < 100; i++ {
		r := RandSeq(6)
		m[r] = m[r] + 1
	}

	for k, v := range m {
		fmt.Println(k, v)
	}
}

func Test_RandomInt(t *testing.T) {
	m := make(map[int]int)
	for i := 0; i < 100; i++ {
		r := RandInt(1, 100)
		m[r] = m[r] + 1
	}

	for k, v := range m {
		fmt.Println(k, v)
	}
}

func Test_RandPerm(t *testing.T) {
	num := 5
	r := RandPerm(num)
	fmt.Println(r)
	if len(r) != num {
		t.Fail()
	}
}
