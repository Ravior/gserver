package gmath

func MinInt(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func MaxInt(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func AbsInt(x int) int {
	if x > 0 {
		return x
	} else {
		return -x
	}
}

func AbsInt32(x int32) int32 {
	if x > 0 {
		return x
	} else {
		return -x
	}
}

// SubAbsInt32 相减-取绝对值
func SubAbsInt32(a int32, b int32) int32 {
	abs := a - b
	if abs < 0 {
		return -abs
	}
	return abs
}

// MaxInt32 最大值
func MaxInt32(a int32, b int32) int32 {
	if a >= b {
		return a
	} else {
		return b
	}
}
