package empty_index

import (
	"fmt"
	"github.com/eosspark/eos-go/common/container/allocator"
	"os"
	"syscall"
	"testing"
	"unsafe"
)

func TestEpIndex_Insert(t *testing.T) {
	alloc = allocator.NewDefaultAllocator(create, nil)
	m := NewEpIndex()
	byID := m.GetByID()

	m.Insert(item{1, 2})

	fmt.Println(byID)

	m.Insert(item{2, 3})

	fmt.Println(byID)

	m.Insert(item{3, 6})

	fmt.Println(byID)
	//fmt.Println(byID.Values())

	itr := byID.Find(2)
	fmt.Println(itr.Key(), itr.Value())

	m.Modify(itr, func(i *item) {
		i.id = 30
		i.num = 2000
	})

	fmt.Println(itr.Key(), itr.Value())
	fmt.Println(byID.Values())

}

func watch(m *EpIndex) {}

func Test_ReadList(t *testing.T) {
	const (
		_TestMemoryFilePath = "/tmp/data/mmap.bin"
		_TestMemorySize     = 1024 * 1024 * 4
	)

	f, _ := os.OpenFile(_TestMemoryFilePath, os.O_RDWR, 0644)
	data, _ := syscall.Mmap(int(f.Fd()), 0, _TestMemorySize, syscall.PROT_READ, syscall.MAP_SHARED)

	//for i := 0; i < 100; i++ {
	//	fmt.Print(data[i], " ")
	//}

	m := (*EpIndex)(unsafe.Pointer(&data[4]))
	byID := m.GetByID()
	fmt.Println(byID.Values())

}

func BenchmarkEpIndex_Insert(b *testing.B) {
	b.StopTimer()

	alloc = allocator.NewDefaultAllocator(create, nil)
	m := NewEpIndex()
	byID := m.GetByID()
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		if itr,ok := byID.Insert(item{uint32(n), uint64(n + 1)}); ok {
			byID.Erase(itr)
		}

	}
}
