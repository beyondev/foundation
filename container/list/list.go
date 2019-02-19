package list

import (
	"foundation/allocator"
	"foundation/offsetptr"
	"unsafe"
)

// template type List(Value,Allocator)
type Value = int

var Allocator allocator.MemoryManager = nil

type Node struct {
	next  offsetptr.Pointer `*Node`
	prev  offsetptr.Pointer `*Node`
	list  offsetptr.Pointer `*List`
	Value Value
}

const _SizeofNode = unsafe.Sizeof(Node{})

func NewNode(value Value) *Node {
	var node *Node
	if Allocator == nil {
		node = new(Node)
	} else {
		node = (*Node)(Allocator.Allocate(_SizeofNode))
	}

	node.next = *offsetptr.NewNil()
	node.prev = *offsetptr.NewNil()
	node.list = *offsetptr.NewNil()
	node.Value = value

	return node
}

func (n *Node) Free() {
	if n != nil && Allocator != nil {
		Allocator.DeAllocate(unsafe.Pointer(n))
	}
}

func (n *Node) Next() *Node {
	if p := (*Node)(n.next.Get()); p != &(*List)(n.list.Get()).root {
		return p
	}
	return nil
}

type List struct {
	root Node
	len  int
}

const _SizeofList = unsafe.Sizeof(List{})

func New() *List {
	if Allocator == nil {
		return new(List).Init()
	} else {
		return (*List)(Allocator.Allocate(_SizeofList)).Init()
	}
}

func (l *List) Free() {
	if l != nil {
		Allocator.DeAllocate(unsafe.Pointer(l))
	}
}

func (l *List) Init() *List {
	l.root.next.Set(unsafe.Pointer(&l.root))
	l.root.prev.Set(unsafe.Pointer(&l.root))
	l.len = 0
	return l
}

func (l *List) lazyInit() {
	if l.len == 0 {
		l.Init()
	}
}

func (l *List) Front() *Node {
	if l.len == 0 {
		return nil
	}

	return (*Node)(l.root.next.Get())
}

func (l *List) Back() *Node {
	if l.len == 0 {
		return nil
	}

	return (*Node)(l.root.prev.Get())
}

func (l *List) insert(e, at *Node) *Node {
	n := offsetptr.NewPointer(at.next.Get())

	at.next.Set(unsafe.Pointer(e))
	e.prev.Set(unsafe.Pointer(at))
	e.next.Set(n.Get())
	(*Node)(n.Get()).prev.Set(unsafe.Pointer(e))
	e.list.Set(unsafe.Pointer(l))

	l.len++
	return e
}

func (l *List) PushFront(value Value) *Node {
	l.lazyInit()
	return l.insert(NewNode(value), &l.root)
}

func (l *List) PushBack(value Value) *Node {
	l.lazyInit()
	return l.insert(NewNode(value), (*Node)(l.root.prev.Get()))
}

func (l *List) Values() []Value {
	if l.len == 0 {
		return nil
	}

	values := make([]Value, 0, l.len)

	for n := l.Front(); n != nil; n = n.Next() {
		values = append(values, n.Value)
	}

	return values
}

func (l *List) remove(n *Node) {
	(*Node)(n.prev.Get()).next.Forward(&n.next)
	(*Node)(n.next.Get()).prev.Forward(&n.prev)
	(*Node)(n.next.Get()).Free()
	(*Node)(n.prev.Get()).Free()
	n.next.Set(nil)
	n.prev.Set(nil)
	n.list.Set(nil)
	l.len--
}

func (l *List) Remove(n *Node) {
	if (*List)(n.list.Get()) == l {
		l.remove(n)
	}
}
