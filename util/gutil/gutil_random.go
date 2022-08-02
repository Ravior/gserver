package gutil

import (
	"bytes"
	"math/rand"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func init() {
	rand.Seed(time.Now().UnixNano()) // 产生随机种子
}

func Random() float64 {
	return RandFloat64Max(1)
}

func RandSeq(size int) string {
	var s bytes.Buffer
	for i := 0; i < size; i++ {
		s.WriteByte(letters[rand.Int63()%int64(len(letters))])
	}
	return s.String()
}

func RandInt(min int, max int) int {
	if min >= max {
		return min
	}
	return rand.Intn(max-min) + min
}

func RandIntMax(max int) int {
	if max <= 0 {
		return 0
	}
	return rand.Intn(max)
}

func RandInt32(min int32, max int32) int32 {
	if min >= max {
		return min
	}
	return rand.Int31n(max-min) + min
}

func RandInt32Max(max int32) int32 {
	if max <= 0 {
		return 0
	}
	return rand.Int31n(max)
}

func RandInt64(min int64, max int64) int64 {
	if min >= max {
		return min
	}
	return rand.Int63n(max-min) + min
}

func RandInt64Max(max int64) int64 {
	if max <= 0 {
		return 0
	}
	return rand.Int63n(max)
}

func RandFloat64(min float64, max float64) float64 {
	if min >= max {
		return min
	}
	return min + rand.Float64()*(max-min)
}

func RandFloat64Max(max float64) float64 {
	if max <= 0 {
		return 0
	}
	return RandFloat64(0, max)
}

// RandPerm returns, as a slice of n int numbers, a pseudo-random permutation of the integers [0,n).
func RandPerm(n int) []int {
	m := make([]int, n)
	for i := 0; i < n; i++ {
		j := RandIntMax(i + 1)
		m[i] = m[j]
		m[j] = i
	}
	return m
}
