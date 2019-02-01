package list

import (
	"fmt"
	"testing"
)

func TestList_Push(t *testing.T) {
	l := New()
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)

	fmt.Println(l.len, l.Values())

	l.PushFront(4)
	l.PushFront(5)
	l.PushFront(6)

	fmt.Println(l.len, l.Values())
}

func TestList_Remove(t *testing.T) {
	l := New()
	n1 := l.PushBack(1)
	n2 := l.PushBack(2)
	n3 := l.PushBack(3)

	fmt.Println(l.Values())

	l.Remove(n1)
	fmt.Println(l.Values())
	l.Remove(n2)
	fmt.Println(l.Values())
	l.Remove(n3)
	fmt.Println(l.Values())


}

func TestNode_Next(t *testing.T) {
	l := New()
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)

	fmt.Println(l.Front().Value)
	fmt.Println(l.Front().Next().Value)
	fmt.Println(l.Front().Next().Next().Value)
	fmt.Println(l.Front().Next().Next().Next())
}