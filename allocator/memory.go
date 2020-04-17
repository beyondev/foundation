package allocator

//#include "string.h"
import "C"

import (
	"errors"
	"unsafe"
)

type MemoryManager interface {
	Allocate(size uintptr) unsafe.Pointer
	DeAllocate(p unsafe.Pointer)
}

type Memory interface {
	Malloc(size uintptr) *byte
	Free(ptr *byte)
}

func Memset(p unsafe.Pointer, ch uint8, len uintptr) {
	C.memset(p, C.int(ch), C.size_t(len))
}

func Memcpy(dest, src unsafe.Pointer, n uintptr) {
	C.memcpy(dest, src, C.size_t(n))
}

var BadAlloc = errors.New("bad alloc")
