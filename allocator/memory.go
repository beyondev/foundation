package allocator

import (
	"unsafe"
)

type MemoryManager interface {
	Allocate(size uintptr) unsafe.Pointer
	DeAllocate(addr unsafe.Pointer)
}

type Memory interface {
	Malloc(size uintptr) *byte
	Free(ptr *byte)
}

type BadAlloc struct {
	Message string
}

func (b BadAlloc) String() string { return b.Message }

var NoAlloc = MemoryManager(nil)

