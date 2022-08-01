package utils

import (
	"fmt"
	"testing"
)

func Test_Trim(t *testing.T) {
	s := " gserver "
	s1 := Trim(s)
	fmt.Println(s1)
	if s1 != "gserver" {
		t.Fail()
	}
}

func Test_IsNumeric(t *testing.T) {
	s := "0"
	if IsNumeric(s) == false {
		t.Fail()
	}
}

func Test_UcFirst(t *testing.T) {
	s := "gserver"
	s1 := UcFirst(s)
	fmt.Println(s1)
	if s1 != "Gserver" {
		t.Fail()
	}
}
