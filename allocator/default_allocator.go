package allocator

import (
	"unsafe"
)

type DefaultAllocator struct {
	availableHeaps []DefaultMemory
	actualHeapSize uint32
	activeHeap     uint32
	activeFreeHeap uint32
	initialized    bool
	grow           func() []byte
}

func NewDefaultAllocator(initialHeap func() []byte, growHeap func() []byte) *DefaultAllocator {
	a := &DefaultAllocator{}
	a.Init(initialHeap, growHeap)
	return a
}

func (a *DefaultAllocator) Init(initialHeap func() []byte, growHeap func() []byte) {
	if initialHeap == nil {
		return
	}
	heap := initialHeap()

	a.availableHeaps = []DefaultMemory{{heap, uintptr(len(heap)), 0}}
	a.actualHeapSize = 1
	a.grow = growHeap
	a.initialized = true
}

func (a *DefaultAllocator) Allocate(size uintptr) unsafe.Pointer {
	if size == 0 {
		return nil
	}

	a.adjustToMemBlock(&size)

	var buf *byte
	var current Memory

	if a.activeHeap < a.actualHeapSize {
		current = &a.availableHeaps[a.activeHeap]
	}

	for current != nil {
		buf = current.Malloc(size)

		if buf != nil {
			break
		}

		current = a.nextActiveHeap()
	}

	if buf == nil {
		endFreeHeap, firstLoop := a.activeFreeHeap, true

		for endFreeHeap != a.activeFreeHeap || firstLoop {
			buf = a.availableHeaps[a.activeFreeHeap].MallocFromFreed(size)

			if buf != nil {
				break
			}

			if a.activeFreeHeap++; a.activeFreeHeap == a.actualHeapSize {
				a.activeFreeHeap = 0
			}

			firstLoop = false
		}
	}

	if buf == nil {
		panic(BadAlloc{"no memory left"})
	}

	//return unsafe.Pointer(buf)
	return unsafe.Pointer(uintptr(unsafe.Pointer(buf)))
}

func (a *DefaultAllocator) DeAllocate(addr unsafe.Pointer) {
	if addr == nil {
		return
	}

	a.availableHeaps[0].Free((*byte)(addr))
}

func (a *DefaultAllocator) nextActiveHeap() Memory {
	if a.grow != nil {
		heap := a.grow()
		m := DefaultMemory{heap, uintptr(len(heap)), 0}
		a.availableHeaps = append(a.availableHeaps, m)
		return &m
	}
	return nil
}

func (a *DefaultAllocator) adjustToMemBlock(size *uintptr) {
	remainder := (*size + _SizeMarker) & _RemMemBlockMask
	if remainder > 0 {
		*size += _MemBlock - remainder
	}
}

type DefaultMemory struct {
	heap     []byte
	heapSize uintptr
	offset   uintptr
}

func (m *DefaultMemory) Malloc(size uintptr) *byte {
	usedUpSize := m.offset + size + _SizeMarker
	if usedUpSize > m.heapSize {
		return nil
	}
	buf := buffer3(&m.heap[m.offset+_SizeMarker], size, &m.heap[m.heapSize-1])
	m.offset += size + _SizeMarker
	buf.markAlloc()
	return buf.ptr
}

func (m *DefaultMemory) Free(ptr *byte) {
	toFree := buffer2(ptr, &m.heap[m.heapSize-1])
	toFree.markFree()
}

func (m *DefaultMemory) MallocFromFreed(size uintptr) *byte {
	//if m.offset != m.heapSize {
	//	panic(BadAlloc{"MallocFromFreed was designed to only be called after heap was completely allocated"})
	//}

	current := _SizeMarker
	for current < m.heapSize {
		currentPtr := &m.heap[current]
		currentBuffer := buffer2(currentPtr, &m.heap[m.heapSize-1])
		if !currentBuffer.isAlloc() {
			if currentBuffer.mergeContiguous(size, false) {
				currentBuffer.markAlloc()
				return currentPtr
			}
		}

		current += currentBuffer.size + _SizeMarker
	}

	return nil
}

func (m *DefaultMemory) isInit() bool {
	return m.heap != nil
}

func (m *DefaultMemory) isInHeap(ptr *byte) bool {
	end := &m.heap[m.heapSize-1]
	first := &m.heap[_SizeMarker]

	return uintptr(unsafe.Pointer(ptr)) >= uintptr(unsafe.Pointer(first)) &&
		uintptr(unsafe.Pointer(ptr)) <= uintptr(unsafe.Pointer(end))
}

type buffer struct {
	ptr     *byte
	size    uintptr
	heapEnd *byte
}

func buffer2(ptr *byte, heapEnd *byte) *buffer {
	return &buffer{ptr, *(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) - _SizeMarker)) & ^_AllocMemoryMask, heapEnd}
}

func buffer3(ptr *byte, size uintptr, heapEnd *byte) *buffer {
	b := &buffer{ptr, size, heapEnd}
	b.setSize(size)
	return b
}

func (b *buffer) setSize(val uintptr) {
	memoryState := *(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(b.ptr)) - _SizeMarker)) & _AllocMemoryMask
	*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(b.ptr)) - _SizeMarker)) = val | memoryState
}

func (b *buffer) markAlloc() {
	*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(b.ptr)) - _SizeMarker)) |= uintptr(_AllocMemoryMask)
}

func (b *buffer) markFree() {
	*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(b.ptr)) - _SizeMarker)) &= ^uintptr(_AllocMemoryMask)
}

func (b *buffer) isAlloc() bool {
	return *(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(b.ptr)) - _SizeMarker))&_AllocMemoryMask > 0
}

func (b *buffer) mergeContiguous(neededSize uintptr, allOrNothing bool) bool {
	if allOrNothing && uintptr(unsafe.Pointer(b.heapEnd))-uintptr(unsafe.Pointer(b.ptr)) < neededSize {
		return false
	}

	possibleSize := b.size
	for possibleSize < neededSize && uintptr(unsafe.Pointer(b.ptr))+possibleSize < uintptr(unsafe.Pointer(b.heapEnd)) {
		nextMemFlagSize := *(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(b.ptr)) + possibleSize))
		if nextMemFlagSize&_AllocMemoryMask > 0 {
			break
		}

		possibleSize += (nextMemFlagSize & ^_AllocMemoryMask) + _SizeMarker
	}

	if allOrNothing && possibleSize < neededSize {
		return false
	}

	// combine
	newSize := neededSize
	if possibleSize < neededSize {
		newSize = possibleSize
	}
	b.setSize(newSize)

	if possibleSize > neededSize {
		freedSize := possibleSize - neededSize - _SizeMarker
		freedRemainder := buffer3((*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(b.ptr))+neededSize+_SizeMarker)), freedSize, b.heapEnd)
		freedRemainder.markFree()
	}

	return newSize == neededSize

}

const _MemBlock = uintptr(8)
const _RemMemBlockMask = _MemBlock - 1

const _SizeMarker = unsafe.Sizeof(uint32(0))
const _AllocMemoryMask = uintptr(1) << 31

//const _SizeMarker = unsafe.Sizeof(uintptr(0))
//const _AllocMemoryMask = uintptr(1) << 63
