package collections

import (
	"testing"

	"github.com/leslie-fei/memcore/gom"
	"github.com/stretchr/testify/assert"
)

func TestRingBuffer(t *testing.T) {
	size := uint64(1024)
	mem := gom.NewMemory(size + uint64(sizeOfRingBuffer))
	_ = mem.Attach()
	defer mem.Detach()

	buffer := (*RingBuffer)(mem.Ptr())
	buffer.Reset(size)

	assert.Equal(t, buffer.IsEmpty(), true)
	assert.Equal(t, buffer.Cap(), size)
	assert.Equal(t, buffer.Len(), uint64(0))
	assert.Equal(t, buffer.IsFull(), false)

	w := make([]byte, 512)
	n, err := buffer.Write(w)
	assert.NoError(t, err)
	assert.Equal(t, len(w), n)
	assert.Equal(t, buffer.IsFull(), false)
	assert.Equal(t, buffer.Len(), uint64(len(w)))
	assert.Equal(t, buffer.IsEmpty(), false)
	assert.Equal(t, buffer.Free(), size-uint64(len(w)))

	r := make([]byte, 512)

	n, err = buffer.Read(r)
	assert.NoError(t, err)
	assert.Equal(t, len(r), n)
	assert.Equal(t, buffer.IsEmpty(), true)
	assert.Equal(t, buffer.IsFull(), false)
	assert.Equal(t, buffer.Free(), size)

	// 写绕回
	w = make([]byte, 1000)
	n, err = buffer.Write(w)
	assert.NoError(t, err)
	assert.Equal(t, len(w), n)
	assert.Equal(t, buffer.IsFull(), false)
	assert.Equal(t, buffer.Len(), uint64(len(w)))
	assert.Equal(t, buffer.IsEmpty(), false)
	assert.Equal(t, buffer.Free(), size-uint64(len(w)))

	// 读绕回
	r = make([]byte, 1000)
	n, err = buffer.Read(r)
	assert.NoError(t, err)
	assert.Equal(t, len(r), n)
	assert.Equal(t, buffer.IsEmpty(), true)
	assert.Equal(t, buffer.IsFull(), false)
	assert.Equal(t, buffer.Free(), size)
}
