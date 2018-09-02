package bitmap

import (
	"errors"
	"fmt"
	"sync"
)

// The Max Size is 0x01 << 32
const bitmapSize = 0x01 << 32

// Bitmap 数据结构定义
type Bitmap struct {
	// 保存实际的 bit 数据
	data []byte
	// 指示该 Bitmap 的 bit 容量
	bitsize uint64
	// 该 Bitmap 被设置为 1 的最大位置（方便遍历）
	maxpos uint64

	mu sync.RWMutex
}

// NewBitmapMax 使用最大容量实例化一个 Bitmap
func NewBitmapMax() *Bitmap {
	bm, _ := NewBitmap(bitmapSize)
	return bm
}

// NewBitmap 根据指定的 size 实例化一个 Bitmap
func NewBitmap(size int) (*Bitmap, error) {
	if size == 0 {
		size = bitmapSize
	} else if remainder := size % 8; remainder != 0 {
		size += 8 - remainder
	}

	if size > bitmapSize {
		return nil, errors.New("size overflows")
	}

	return &Bitmap{data: make([]byte, size>>3), bitsize: uint64(size - 1)}, nil
}

// Add 设 offset 位置的 bit 置为1
func (bm *Bitmap) Add(offset uint64) bool {
	index, pos := offset/8, offset%8
	if bm.bitsize < offset {
		return false
	}

	bm.mu.Lock()
	bm.data[index] |= 0x01 << pos
	bm.mu.Unlock()

	if bm.maxpos < offset {
		bm.maxpos = offset
	}

	return true
}

// Del 设 offset 位置的 bit 为0
func (bm *Bitmap) Del(offset uint64) bool {
	index, pos := offset/8, offset%8
	if bm.bitsize < offset {
		return false
	}

	bm.mu.Lock()
	bm.data[index] &^= 0x01 << pos
	bm.mu.Unlock()

	if bm.maxpos <= offset {
		bm.maxpos, _ = bm.Prev(offset)
	}

	return true
}

// Prev ...
func (bm *Bitmap) Prev(offset uint64) (uint64, bool) {
	var i uint64
	for i = offset; i >= 0; i-- {
		if bm.Has(i) {
			return i, true
		}
	}
	return 0, false
}

// Next ...
func (bm *Bitmap) Next(offset uint64) (uint64, bool) {
	var i uint64
	for i = offset; i <= bm.bitsize; i++ {
		if bm.Has(i) {
			return i, true
		}
	}
	return 0, false
}

// Has 获得 offset 位置处的 value
func (bm *Bitmap) Has(offset uint64) bool {
	if bm.bitsize < offset {
		return false
	}
	index, pos := offset/8, offset%8

	bm.mu.RLock()
	has := (bm.data[index]>>pos)&0x01 == 1
	bm.mu.RUnlock()

	return has
}

// Maxpos 获取 bitmap 存在的最大位置
func (bm *Bitmap) Maxpos() uint64 {
	return bm.maxpos
}

// BitSize ...
func (bm *Bitmap) BitSize() uint64 {
	return bm.bitsize
}

// String 实现 Stringer 接口（只输出开始的100个元素）
func (bm *Bitmap) String() string {
	var (
		maxTotal, bitTotal uint64 = 100, bm.maxpos + 1
		offset             uint64
	)

	if bm.maxpos > maxTotal {
		bitTotal = maxTotal
	}

	bm.mu.RLock()
	defer bm.mu.RUnlock()

	numSlice := make([]uint64, 0, bitTotal)
	for offset = 0; offset < bitTotal; offset++ {
		if bm.Has(offset) {
			numSlice = append(numSlice, offset)
		}
	}

	return fmt.Sprintf("%v", numSlice)
}
