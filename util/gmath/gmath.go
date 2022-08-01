package gmath

func MinInt(x, y int) int {
	if x > y {
		return y
	}
	return x
}

// MinInt32 最大值
func MinInt32(x, y int32) int32 {
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

// MaxInt32 最大值
func MaxInt32(x, y int32) int32 {
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

// SubAbsInt 相减-取绝对值
func SubAbsInt(x, y int) int {
	abs := x - y
	if abs < 0 {
		return -abs
	}
	return abs
}

// SubAbsInt32 相减-取绝对值
func SubAbsInt32(x, y int32) int32 {
	abs := x - y
	if abs < 0 {
		return -abs
	}
	return abs
}
