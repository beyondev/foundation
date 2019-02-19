package hashtable

import (
	"bytes"
	"fmt"
	"foundation/allocator"
	"foundation/allocator/callocator"
	"foundation/container"
	. "foundation/offsetptr"
	"unsafe"
)

type K = int
type V = int

var hashfunc = container.Hash

var Allocator allocator.MemoryManager = callocator.Instance

type Hashtable struct {
	b Pointer `*Buckets`
	l uintptr
}

const _SizeofHashtable = unsafe.Sizeof(Hashtable{})

func New(size int) *Hashtable {
	h := (*Hashtable)(Allocator.Allocate(_SizeofHashtable))
	h.b.Set(unsafe.Pointer(NewBuckets(uintptr(size))))
	h.l = 0
	return h
}

func (h *Hashtable) Len() int {
	return int(h.l)
}

func (h *Hashtable) Cap() int {
	return int((*Buckets)(h.b.Get()).len)
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
		(*Buckets)(h.b.Get()).Free()
		Allocator.DeAllocate(unsafe.Pointer(h))
	}
}

func (h *Hashtable) Put(key K, value V) {
	b := (*Buckets)(h.b.Get())
	bucket := hashfunc(key) % b.len
	entry, appended := b.At(bucket).Put(bucket, key, value)

	b.Put(bucket, entry)

	if appended {
		entry.bucket = bucket
		h.l++
	}

	if h.l > b.len {
		h.resize()
	}
}

func (h *Hashtable) Remove(key K) {
	b := (*Buckets)(h.b.Get())
	bucket := hashfunc(key) % b.len
	if entry, removed := b.At(bucket).Remove(key); removed {
		b.Put(bucket, entry)
		h.l--
	}
}

func (h *Hashtable) Get(key K) (V, bool) {
	b := (*Buckets)(h.b.Get())
	return b.At(hashfunc(key) % b.len).Get(key)
}

func (h *Hashtable) Each(f func(key K, value V)) {
	b := (*Buckets)(h.b.Get())
	for i := uintptr(0); i < b.len; i++ {
		for e := b.At(i); e != nil; e = (*Entry)(e.next.Get()) {
			f(e.key, e.value)
		}
	}
}

func (h *Hashtable) resize() {
	b := (*Buckets)(h.b.Get())
	n := b.len * 2
	tmp := NewBuckets(n)
	for bucket := uintptr(0); bucket < b.len; bucket++ {
		first := b.At(bucket)
		for first != nil {
			newBucket := hashfunc(first.key) % n
			b.Put(bucket, (*Entry)(first.next.Get()))
			first.bucket = newBucket
			first.next.Set(unsafe.Pointer(tmp.At(newBucket)))
			tmp.Put(newBucket, first)
			first = b.At(bucket)
		}
	}

	Allocator.DeAllocate(h.b.Get())
	h.b.Set(unsafe.Pointer(tmp))
}

type Buckets struct {
	array Pointer `Pointer<*Entry>`
	len   uintptr
}

const _SizeofBuckets = unsafe.Sizeof(Buckets{})
const _SizeofBucket = unsafe.Sizeof(Pointer{})

func NewBuckets(size uintptr) *Buckets {
	b := (*Buckets)(Allocator.Allocate(_SizeofBuckets))
	if size == 0 {
		size++
	}
	b.array.Set(Allocator.Allocate(_SizeofBucket * size))
	allocator.Memset(b.array.Get(), 0, _SizeofBucket*size)
	b.len = size
	return b
}

func (b *Buckets) Free() {
	if b != nil {
		Allocator.DeAllocate(b.array.Get())
		Allocator.DeAllocate(unsafe.Pointer(b))
	}
}

func (b *Buckets) Put(i uintptr, e *Entry) {
	(*Pointer)(unsafe.Pointer(uintptr(b.array.Get()) + _SizeofBucket*i)).Set(unsafe.Pointer(e))
	//*(**Entry)(unsafe.Pointer(uintptr(unsafe.Pointer(b.array)) + _SizeofBucket*i)) = e
}
func (b *Buckets) At(i uintptr) *Entry {
	p := (*Pointer)(unsafe.Pointer(uintptr(b.array.Get()) + _SizeofBucket*i))
	if p.IsSelf() {
		p.Set(nil)
	}
	//fmt.Println("addr of **b.array:", uintptr(p.Get()))

	return (*Entry)(p.Get())
	//return *(**Entry)(unsafe.Pointer(uintptr(unsafe.Pointer(b.array)) + _SizeofBucket*i))
}

type Entry struct {
	key    K
	value  V
	bucket uintptr
	next   Pointer `*Entry`
}

const _SizeofEntry = unsafe.Sizeof(Entry{})

func NewEntry(bucket uintptr, key K, value V) *Entry {
	e := (*Entry)(Allocator.Allocate(_SizeofEntry))
	e.bucket = bucket
	e.key = key
	e.value = value
	e.next.Set(nil)
	return e
}

func (e *Entry) Free() {
	if e != nil {
		e.next.Set(nil)
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
		enext, appended := (*Entry)(e.next.Get()).Put(bucket, key, value)
		e.next.Set(unsafe.Pointer(enext))
		return e, appended
	}
}

func (e *Entry) Remove(key K) (entry *Entry, removed bool) {
	if e == nil {
		return e, false
	}

	if e.key == key {
		next := (*Entry)(e.next.Get())
		e.Free()
		return next, true
	} else {
		enext, removed := (*Entry)(e.next.Get()).Remove(key)
		e.next.Set(unsafe.Pointer(enext))
		return e, removed
	}
}

func (e *Entry) Get(key K) (value V, has bool) {
	for entry := e; entry != nil; entry = (*Entry)(entry.next.Get()) {
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

	b := (*Buckets)(h.b.Get())
	for i := uintptr(0); i < b.len; i++ {
		if e := b.At(i); e != nil {
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
	if !itr.e.next.IsNil() {
		itr.e = (*Entry)(itr.e.next.Get())
		return true
	}

	b := (*Buckets)(itr.h.b.Get())
	for bucket := itr.e.bucket + 1; bucket < b.len; bucket++ {
		if entry := b.At(bucket); entry != nil {
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
