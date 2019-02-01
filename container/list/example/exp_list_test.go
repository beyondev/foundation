package example

import (
	"container/list"
	"fmt"
	"os"
	"syscall"
	"testing"
	"unsafe"
)

func Test_insert(t *testing.T) {
	l := NewExampleList()
	//p := byte(1)
	l.PushBack(Item{1, [4]byte{'a', 'b'}})
	l.PushBack(Item{2, [4]byte{'a', 'b', 'c'}})
	l.PushBack(Item{3, [4]byte{'a', 'b', 'c', 'd'}})
}

func Test_ReadList(t *testing.T) {
	f, _ := os.OpenFile(_TestMemoryFilePath, os.O_RDWR, 0644)
	data, _ := syscall.Mmap(int(f.Fd()), 0, _TestMemorySize, syscall.PROT_READ, syscall.MAP_SHARED)

	l := (*ExampleList)(unsafe.Pointer(&data[4]))
	for _, item := range l.Values() {
		fmt.Println(item.id, string(item.name[:]))
	}
}

func BenchmarkExampleList_PushBack(b *testing.B) {
	b.StopTimer()
	l := NewExampleList()
	//b.N = 1000000
	b.StartTimer()
	for n := uint32(0); n < uint32(b.N); n++ {
		node := l.PushBack(Item{n, [4]byte{'a', 'b', 'c', 'd'}})
		l.Remove(node)
	}
}

func BenchmarkList_PushBack(b *testing.B) {
	b.StopTimer()
	l := list.New()
	b.N = 1000000
	b.StartTimer()
	for n := uint32(0); n < uint32(b.N); n++ {
		l.PushBack(Item{n, [4]byte{'a', 'b', 'c', 'd'}})
	}
}
