package slice

import (
	"bytes"
	"fmt"
	"unsafe"

	"github.com/Beyond-simplechain/foundation/allocator"
	"github.com/Beyond-simplechain/foundation/allocator/callocator"
	"github.com/Beyond-simplechain/foundation/container"
)

type V = int

var Allocator allocator.MemoryManager = callocator.Instance

type Slice struct {
	array unsafe.Pointer
	len   uintptr
	cap   uintptr
}

const _SizeofSlice = unsafe.Sizeof(Slice{})
const _SizeofV = unsafe.Sizeof(*new(V))

func New(cap int) *Slice {
	s := (*Slice)(Allocator.Allocate(_SizeofSlice))
	s.cap = uintptr(cap)
	s.len = 0
	s.array = nil

	if cap > 0 {
		s.array = Allocator.Allocate(uintptr(cap) * _SizeofV)
	}

	return s

}

func (s *Slice) Free() {
	if s != nil {
		Allocator.DeAllocate(unsafe.Pointer(s))
	}
}

func (s *Slice) Len() int {
	return int(s.len)
}

func (s *Slice) Cap() int {
	return int(s.cap)
}

func (s *Slice) String() string {
	var buffer bytes.Buffer
	buffer.WriteByte('[')
	for i := 0; i < s.Len(); i++ {
		if i != 0 {
			buffer.WriteByte(' ')
		}
		buffer.WriteString(fmt.Sprint(s.Get(i)))
	}
	buffer.WriteByte(']')
	return buffer.String()
}

func (s *Slice) PushBack(value V) {
	if s.array == nil {
		s.array = Allocator.Allocate(_SizeofV)
		s.cap = 1
	}

	if s.len+1 > s.cap {
		s.expend()
	}

	*(*V)(unsafe.Pointer(uintptr(s.array) + s.len*_SizeofV)) = value

	s.len++
}

func (s *Slice) Put(i int, value V) {
	if uintptr(i) >= s.len {
		panic(container.ErrIndexOutOfRange)
	}

	*(*V)(unsafe.Pointer(uintptr(s.array) + uintptr(i)*_SizeofV)) = value
}

func (s *Slice) Get(i int) V {
	if uintptr(i) >= s.len {
		panic(container.ErrIndexOutOfRange)
	}

	return *(*V)(unsafe.Pointer(uintptr(s.array) + uintptr(i)*_SizeofV))
}

func (s *Slice) expend() {
	src := s.array
	s.cap <<= 1
	s.array = Allocator.Allocate(_SizeofV * s.cap)
	allocator.Memcpy(s.array, src, _SizeofV*s.len)
	Allocator.DeAllocate(unsafe.Pointer(src))
}
