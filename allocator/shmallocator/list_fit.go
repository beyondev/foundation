//go:build none
// +build none

package shmallocator

import (
	"unsafe"

	"github.com/beyondev/foundation/allocator"
)

type SimpleSeqFit struct {
	mHead *ListFitNode
	mHeap []byte
}

type blockCtrl struct {
	next *blockCtrl
	size uintptr
}

type header struct {
	mRoot          blockCtrl
	mAllocated     uintptr
	mSize          uintptr
	mExtraHdrBytes uintptr
}

type ListFitNode struct {
	prev *ListFitNode
	next *ListFitNode
	size uintptr
	free bool
}

const _SizeofListFitNode = unsafe.Sizeof(ListFitNode{})

func NewListFit() *ListFit {
	l := new(ListFit)
	l.heap = make([]byte, 1024*1024)
	l.head = (*ListFitNode)(unsafe.Pointer(&l.heap[0]))
	l.head.next = nil
	l.head.prev = nil
	l.head.size = 1024 * 1024
	l.head.free = true
	l.tail = l.head
	return l
}

func (l *ListFit) Malloc(size uintptr) unsafe.Pointer {
	var node, ptr = l.head, unsafe.Pointer(nil)

	for node.free && node != nil {
		if node.size >= size {
			//TODO get and divide
		} else {
			node = node.next
		}
	}

	if ptr == nil {
		panic(allocator.BadAlloc)
	}

	return ptr
}

func (l *ListFit) Free(ptr unsafe.Pointer) {
	node := (*ListFitNode)(unsafe.Pointer(uintptr(ptr) - _SizeofListFitNode))
	node.free = true
	//	TODO
}
