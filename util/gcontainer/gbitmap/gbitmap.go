package gbitmap

import (
	"math"
	"unsafe"
)

const (
	sizeInt = int(unsafe.Sizeof(int(0))) * 8
)

type (
	BitMap struct {
		m_Bits []int
		m_Size int
	}

	IBitMap interface {
		Init(size int)
		Set(index int)       //设置位
		Test(index int) bool //位是否被设置
		Clear(index int)     //清楚位
		ClearAll()
	}
)

func (b *BitMap) Init(size int) {
	b.m_Size = int(math.Ceil(float64(size) / float64(sizeInt)))
	b.m_Bits = make([]int, b.m_Size)
}

func (b *BitMap) Set(index int) {
	if index >= b.m_Size*sizeInt {
		return
	}

	b.m_Bits[index/sizeInt] |= 1 << uint(index%sizeInt)
}

func (b *BitMap) Test(index int) bool {
	if index >= b.m_Size*sizeInt {
		return false
	}

	return b.m_Bits[index/sizeInt]&(1<<uint(index%sizeInt)) != 0

}

func (b *BitMap) Clear(index int) {
	if index >= b.m_Size*sizeInt {
		return
	}

	b.m_Bits[index/sizeInt] &= ^(1 << uint(index%sizeInt))
}

func (b *BitMap) ClearAll() {
	b.Init(b.m_Size * sizeInt)
}

func NewBitMap(size int) *BitMap {
	bitmap := &BitMap{}
	bitmap.Init(size)
	return bitmap
}
