package collections

import (
	"fmt"
	"unsafe"

	"github.com/leslie-fei/memcore"
)

var sizeOfRingBuffer = unsafe.Sizeof(RingBuffer{})

type RingBuffer struct {
	size     uint64
	readIdx  uint64
	writeIdx uint64
	isFull   bool
}

func (r *RingBuffer) Reset(size uint64) {
	*r = RingBuffer{}
	r.size = size
}

func (r *RingBuffer) Len() uint64 {
	// writeIdx == readIdx, 1: empty 2 : full
	if r.writeIdx == r.readIdx {
		if r.isFull {
			return r.size
		}
		return 0
	}
	if r.writeIdx < r.readIdx {
		return r.size - r.readIdx + r.writeIdx
	}
	return r.writeIdx - r.readIdx
}

func (r *RingBuffer) Cap() uint64 {
	return r.size
}

func (r *RingBuffer) IsFull() bool {
	return r.isFull
}

func (r *RingBuffer) IsEmpty() bool {
	return r.Len() == 0
}

func (r *RingBuffer) Read(b []byte) (n int, err error) {
	if len(b) == 0 || r.IsEmpty() {
		return 0, nil
	}
	ptr := unsafe.Pointer(&b)
	var length = r.Len()
	length = min(length, uint64(len(b)))
	// 还没写入绕回
	if r.readIdx < r.writeIdx {
		memcore.Memmove(ptr, r.bufferPtr(r.readIdx), uintptr(length))
	} else {
		// 已经写入绕回
		// 绕回之前可读的长度
		right := r.size - r.readIdx
		// 如果绕回前足够
		if right >= length {
			memcore.Memmove(ptr, r.bufferPtr(r.readIdx), uintptr(length))
		} else {
			// 分段读取
			memcore.Memmove(ptr, r.bufferPtr(r.readIdx), uintptr(right))
			memcore.Memmove(unsafe.Pointer(uintptr(ptr)+uintptr(right)), r.bufferPtr(0), uintptr(length-right))
		}
	}

	r.readIdx = (r.readIdx + length) % r.size
	r.isFull = false

	return int(length), nil
}

func (r *RingBuffer) Write(b []byte) (n int, err error) {
	if len(b) == 0 || r.IsFull() {
		return 0, nil
	}

	writeable := r.Free()
	writeable = min(writeable, uint64(len(b)))

	ptr := unsafe.Pointer(&b)
	// 已经绕回, 或者没绕回还有足够空间, 不需要分段写
	if r.writeIdx > r.readIdx || r.size-r.writeIdx >= writeable {
		memcore.Memmove(r.bufferPtr(r.writeIdx), ptr, uintptr(writeable))
	} else {
		// 分段写
		right := r.size - r.writeIdx
		left := writeable - right
		memcore.Memmove(r.bufferPtr(r.writeIdx), ptr, uintptr(right))
		memcore.Memmove(r.bufferPtr(0), unsafe.Pointer(uintptr(ptr)+uintptr(right)), uintptr(left))
	}

	r.writeIdx = (r.writeIdx + writeable) % r.size
	if r.writeIdx == r.readIdx {
		r.isFull = true
	}

	return int(writeable), nil
}

func (r *RingBuffer) Free() uint64 {
	return r.size - r.Len()
}

func (r *RingBuffer) bufferPtr(offset uint64) unsafe.Pointer {
	if offset > r.size {
		panic(fmt.Errorf("offset out of range, offset: %d, size: %d", offset, r.size))
	}
	return unsafe.Pointer(uintptr(unsafe.Pointer(r)) + sizeOfRingBuffer + uintptr(offset))
}
