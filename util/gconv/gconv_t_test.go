package gconv

import "testing"

func Test_Int32(t *testing.T) {
	if Int32("10001") != 10001 {
		t.Fail()
	}
}

func Test_Int64(t *testing.T) {
	if Int64("10001") != 10001 {
		t.Fail()
	}
}

func Test_Float32(t *testing.T) {
	if Float32("10001.00") != 10001.00 {
		t.Fail()
	}
}

func Test_Float64(t *testing.T) {
	if Float64("10001.00") != 10001.00 {
		t.Fail()
	}
}

func Test_String(t *testing.T) {
	if String(10001) != "10001" {
		t.Fail()
	}
}

func Test_IntToString(t *testing.T) {
	if IntToString(10001) != "10001" {
		t.Fail()
	}
}

func Test_Int32ToString(t *testing.T) {
	if Int32ToString(10001) != "10001" {
		t.Fail()
	}
}

func Test_Bool(t *testing.T) {
	if Bool("1") != true {
		t.Fail()
	}

	if Bool("0") != false {
		t.Fail()
	}
}
