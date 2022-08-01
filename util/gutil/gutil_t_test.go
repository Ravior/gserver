package gutil

import (
	"fmt"
	"testing"
)

func Test_Try(t *testing.T) {
	s := `gutil Try test`
	err := Try(func() {
		panic(s)
	})
	if err.Error() != s {
		t.Fail()
	}
}

func Test_TryCatch(t *testing.T) {
	s := `gutil TryCatch test`
	TryCatch(func() {
		panic(s)
	}, func(exception error) {
		if exception.Error() != s {
			t.Fail()
		}
	})
}

func Test_NickCall(t *testing.T) {
	s := `gutil NickCall test`
	NiceCallFunc(func() {
		panic(s)
	})
}

func Test_IsEmpty(t *testing.T) {
	s := ""
	if IsEmpty(s) == false {
		t.Fail()
	}

	s1 := "hello world"
	if IsEmpty(s1) == true {
		t.Fail()
	}

	s2 := 0
	if IsEmpty(s2) == false {
		t.Fail()
	}

	s3 := make([]int32, 0)
	if IsEmpty(s3) == false {
		t.Fail()
	}

	s4 := make(map[int32]int32, 0)
	if IsEmpty(s4) == false {
		t.Fail()
	}
}

func Test_Keys(t *testing.T) {
	s := map[int32]string{
		1: "first",
		2: "second",
	}
	// [1 2]
	fmt.Println(Keys(s))
	// [first second]
	fmt.Println(Values(s))
	if len(Keys(s)) != len(s) {
		t.Fail()
	}
}
