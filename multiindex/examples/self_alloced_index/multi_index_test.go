package self_alloced_index

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"syscall"
	"testing"
	"unsafe"
)

func TestIndex_Insert(t *testing.T) {
	m := NewTestIndex()
	m.Insert(SharedClass{1, 1})
}

func TestIndex_CRUD(t *testing.T) {
	m := NewTestIndex()
	for n := 0; n < 5; n++ {
		m.Insert(SharedClass{2 * n, uint32(n)})
	}

	f, _ := os.OpenFile(_TestMemoryFilePath, os.O_RDWR|os.O_CREATE, 0644)
	bys := make([]byte, 200)
	f.Read(bys)
	fmt.Println(bys)
	data, _ := syscall.Mmap(int(f.Fd()), 0, _TestMemorySize, syscall.PROT_WRITE, syscall.MAP_SHARED|syscall.MAP_COPY)
	f.Close()

	idx := (*TestIndex)(unsafe.Pointer(&data[4]))
	byId := idx.super
	fmt.Println(byId.Values())

	//itr := byId.Find(6)
	//byId.Modify(itr, func(value *SharedClass) {
	//	value.id = 100
	//	value.num = 999
	//})
}

func TestRecover(t *testing.T) {
	f, err := os.OpenFile(_TestMemoryFilePath, os.O_RDWR|os.O_CREATE, 0644)
	assert.NoError(t, err)
	bys := make([]byte, 200)
	f.Read(bys)
	fmt.Println(bys)
	data, err := syscall.Mmap(int(f.Fd()), 0, _TestMemorySize, syscall.PROT_WRITE, syscall.MAP_SHARED|syscall.MAP_COPY)
	assert.NoError(t, err)
	f.Close()
	defer os.Remove(_TestMemoryFilePath)

	idx := (*TestIndex)(unsafe.Pointer(&data[4]))
	byId := idx.super
	fmt.Println(byId.Values())
}

func BenchmarkIndex_Insert(b *testing.B) {
	m := NewTestIndex()
	b.N = 1000000
	for n := 0; n < b.N; n++ {
		m.Insert(SharedClass{2 * n, uint32(n)})
	}
}
