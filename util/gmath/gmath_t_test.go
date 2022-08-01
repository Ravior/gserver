package gmath

import "testing"

func Test_MinInt(t *testing.T) {
	if MinInt(1, 2) != 1 {
		t.Fail()
	}
}

func Test_MinInt32(t *testing.T) {
	if MinInt32(1, 2) != 1 {
		t.Fail()
	}
}

func Test_MaxInt(t *testing.T) {
	if MaxInt(1, 2) != 2 {
		t.Fail()
	}
}

func Test_MaxInt32(t *testing.T) {
	if MaxInt32(1, 2) != 2 {
		t.Fail()
	}
}

func Test_AbsInt(t *testing.T) {
	if AbsInt(-1) != 1 {
		t.Fail()
	}
}

func Test_AbsInt32(t *testing.T) {
	if AbsInt32(-1) != 1 {
		t.Fail()
	}
}

func Test_SubAbsInt(t *testing.T) {
	if SubAbsInt(1, 3) != 2 {
		t.Fail()
	}
}

func Test_SubAbsInt32(t *testing.T) {
	if SubAbsInt32(1, 3) != 2 {
		t.Fail()
	}
}
