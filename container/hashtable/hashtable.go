//go:build .
// +build .

package hashtable

import (
	"bytes"
	"fmt"
	"unsafe"

	"github.com/beyondev/foundation/allocator"
	"github.com/beyondev/foundation/allocator/callocator"
	"github.com/beyondev/foundation/container"
)

type K = int
type V = int

var hashfunc = container.Hash

var Allocator allocator.MemoryManager = callocator.Instance

type Hashtable struct {
	b *Buckets
	l uintptr
}

const _SizeofHashtable = unsafe.Sizeof(Hashtable{})

func New(size int) *Hashtable {
	h := (*Hashtable)(Allocator.Allocate(_SizeofHashtable))
	h.b = NewBuckets(uintptr(size))
	h.l = 0
	return h
}

func (h *Hashtable) Len() int {
	return int(h.l)
}

func (h *Hashtable) Cap() int {
	return int(h.b.len)
}

func (h *Hashtable) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("hashtable[")
	for first, i := true, h.Begin(); !i.IsEnd(); i.Next() {
		if first {
			first = false
		} else {
			buffer.WriteByte(' ')
		}
		buffer.WriteString(fmt.Sprint(i.Key()))
		buffer.WriteByte(':')
		buffer.WriteString(fmt.Sprint(i.Value()))
	}
	buffer.WriteByte(']')

	return buffer.String()
}

func (h *Hashtable) Free() {
	if h != nil {
		h.b.Free()
		Allocator.DeAllocate(unsafe.Pointer(h))
	}
}

func (h *Hashtable) Put(key K, value V) {
	bucket := hashfunc(key) % h.b.len
	entry, appended := h.b.At(bucket).Put(bucket, key, value)

	h.b.Put(bucket, entry)

	if appended {
		entry.bucket = bucket
		h.l++
	}

	if h.l > h.b.len {
		h.resize()
	}
}

func (h *Hashtable) Remove(key K) {
	bucket := hashfunc(key) % h.b.len
	if entry, removed := h.b.At(bucket).Remove(key); removed {
		h.b.Put(bucket, entry)
		h.l--
	}
}

func (h *Hashtable) Get(key K) (V, bool) {
	return h.b.At(hashfunc(key) % h.b.len).Get(key)
}

func (h *Hashtable) Each(f func(key K, value V)) {
	for i := uintptr(0); i < h.b.len; i++ {
		for e := h.b.At(i); e != nil; e = e.next {
			f(e.key, e.value)
		}
	}
}

func (h *Hashtable) resize() {
	n := h.b.len * 2
	tmp := NewBuckets(n)
	for bucket := uintptr(0); bucket < h.b.len; bucket++ {
		first := h.b.At(bucket)
		for first != nil {
			newBucket := hashfunc(first.key) % n
			h.b.Put(bucket, first.next)
			first.bucket = newBucket
			first.next = tmp.At(newBucket)
			tmp.Put(newBucket, first)
			first = h.b.At(bucket)
		}
	}

	h.b.Free()
	h.b = tmp
}

type Buckets struct {
	array **Entry `Pointer<*Entry>`
	len   uintptr
}

const _SizeofBuckets = unsafe.Sizeof(Buckets{})
const _SizeofBucket = unsafe.Sizeof(&Entry{})

func NewBuckets(size uintptr) *Buckets {
	b := (*Buckets)(Allocator.Allocate(_SizeofBuckets))
	if size == 0 {
		size++
	}
	b.array = (**Entry)(Allocator.Allocate(_SizeofBucket * size))
	allocator.Memset(unsafe.Pointer(b.array), 0, _SizeofBucket*size)
	b.len = size
	return b
}

func (b *Buckets) Free() {
	if b != nil {
		Allocator.DeAllocate(unsafe.Pointer(b.array))
		Allocator.DeAllocate(unsafe.Pointer(b))
	}
}

func (b *Buckets) Put(i uintptr, e *Entry) {
	*(**Entry)(unsafe.Pointer(uintptr(unsafe.Pointer(b.array)) + _SizeofBucket*i)) = e
}

func (b *Buckets) At(i uintptr) *Entry {
	return *(**Entry)(unsafe.Pointer(uintptr(unsafe.Pointer(b.array)) + _SizeofBucket*i))
}

type Entry struct {
	key    K
	value  V
	bucket uintptr
	next   *Entry `*Entry`
}

const _SizeofEntry = unsafe.Sizeof(Entry{})

func NewEntry(bucket uintptr, key K, value V) *Entry {
	e := (*Entry)(Allocator.Allocate(_SizeofEntry))
	e.bucket = bucket
	e.key = key
	e.value = value
	e.next = nil
	return e
}

func (e *Entry) Free() {
	if e != nil {
		Allocator.DeAllocate(unsafe.Pointer(e))
	}
}

func (e *Entry) Put(bucket uintptr, key K, value V) (entry *Entry, appended bool) {
	if e == nil {
		return NewEntry(bucket, key, value), true
	}

	if e.key == key {
		e.value = value
		return e, false
	} else {
		e.next, appended = e.next.Put(bucket, key, value)
		return e, appended
	}
}

func (e *Entry) Remove(key K) (entry *Entry, removed bool) {
	if e == nil {
		return e, false
	}

	if e.key == key {
		next := e.next
		e.Free()
		return next, true
	} else {
		e.next, removed = e.next.Remove(key)
		return e, removed
	}
}

func (e *Entry) Get(key K) (value V, has bool) {
	for entry := e; entry != nil; entry = entry.next {
		if entry.key == key {
			return entry.value, true
		}
	}
	return
}

type Iterator struct {
	h *Hashtable
	e *Entry
}

func (h *Hashtable) Begin() Iterator {
	if h.l == 0 {
		return h.End()
	}

	for i := uintptr(0); i < h.b.len; i++ {
		if e := h.b.At(i); e != nil {
			return Iterator{h, e}
		}
	}

	panic(container.ErrFatalAddress)
}

func (h *Hashtable) End() Iterator {
	return Iterator{h, nil}
}

func (itr Iterator) IsEnd() bool {
	return itr.e == nil
}

func (itr *Iterator) Next() bool {
	if itr.e.next != nil {
		itr.e = itr.e.next
		return true
	}

	for bucket := itr.e.bucket + 1; bucket < itr.h.b.len; bucket++ {
		if entry := itr.h.b.At(bucket); entry != nil {
			itr.e = entry
			return true
		}
	}

	itr.e = nil
	return false
}

func (itr *Iterator) Key() K {
	return itr.e.key
}

func (itr *Iterator) Value() V {
	return itr.e.value
}
