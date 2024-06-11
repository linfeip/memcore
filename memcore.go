package memcore

import "unsafe"

//go:linkname memmove runtime.memmove
//go:noescape
func memmove(dst, src unsafe.Pointer, size uintptr)

//go:linkname memequal runtime.memequal
//go:noescape
func memequal(a, b unsafe.Pointer, size uintptr) bool

func Memmove(dst, src unsafe.Pointer, size uintptr) {
	memmove(dst, src, size)
}

func Memequal(a, b unsafe.Pointer, size uintptr) bool {
	return memequal(a, b, size)
}
