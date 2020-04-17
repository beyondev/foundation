package shmallocator

import (
	"fmt"
	"log"
	"os"
	"syscall"
	"testing"
	"unsafe"

	"github.com/Beyond-simplechain/foundation/offsetptr"
	"github.com/stretchr/testify/assert"
)

const _TestMemorySize = 1024 * 1024 * 1024
const _TestMemoryFilePath = "/tmp/data/mmap.bin"

func create() []byte {
	f, err := syscall.Open(_TestMemoryFilePath, os.O_RDWR|os.O_CREATE, 0644)

	if nil != err {
		log.Fatalln(err)
	}

	// extend file
	if _, err := syscall.Write(f, make([]byte, _TestMemorySize)); nil != err {
		log.Fatalln(err)
	}
	data, err := syscall.Mmap(f, 0, _TestMemorySize, syscall.PROT_WRITE, syscall.MAP_SHARED|syscall.MAP_COPY)

	//data, err := unix.Mmap(int(f.Fd()), 0, 1<<8, syscall.PROT_WRITE, syscall.MAP_SHARED|syscall.MAP_COPY)
	if nil != err {
		log.Fatalln(err)
	}

	if err := syscall.Close(f); nil != err {
		log.Fatalln(err)
	}

	return data
}
func open() []byte {
	f, err := syscall.Open(_TestMemoryFilePath, os.O_RDWR|os.O_CREATE, 0644)

	if nil != err {
		log.Fatalln(err)
	}

	data, err := syscall.Mmap(f, 0, _TestMemorySize, syscall.PROT_WRITE, syscall.MAP_SHARED|syscall.MAP_COPY)

	//data, err := unix.Mmap(int(f.Fd()), 0, 1<<8, syscall.PROT_WRITE, syscall.MAP_SHARED|syscall.MAP_COPY)
	if nil != err {
		log.Fatalln(err)
	}

	if err := syscall.Close(f); nil != err {
		log.Fatalln(err)
	}

	return data
}

func destroy() {
	os.Remove(_TestMemoryFilePath)
}

func close(d []byte) {
	syscall.Munmap(d)
}

var alloc *Allocator

type st struct {
	a uint64
	b uint64
	c offsetptr.Pointer `*uint32`
}

func (s *st) GetC() *uint32 {
	return (*uint32)(s.c.Get())
}

func TestDefaultAllocator_Allocate(t *testing.T) {
	alloc = New(create, nil)

	sta := (*st)(alloc.Allocate(unsafe.Sizeof(st{})))
	sta.a = 1<<16 - 1
	sta.c.Set(alloc.Allocate(unsafe.Sizeof(uint32(0))))
	*sta.GetC() = 255

	//sta.b = 1<<16 - 1
	//sta.c = (*uint32)(alloc.Allocate(unsafe.Sizeof(uint32(0))))
	//sta.SetC(10)
	//
	fmt.Println(sta)
	assert.NoError(t, syscall.Munmap(alloc.availableHeaps[0].heap))

	TestReadMMapData(t)
	//for i:=0; i<100; i++ {
	//	fmt.Print(alloc.availableHeaps[0].heap[i], " ")
	//}

}

func TestReadMMapData(t *testing.T) {
	alloc := New(open, nil)

	sta := (*st)(unsafe.Pointer(&alloc.availableHeaps[0].heap[_SizeMarker]))

	//
	fmt.Println(sta.a)
	fmt.Println(sta.b)
	fmt.Println(*sta.GetC())
	//fmt.Printf("addr sta: %p\n", sta)
	//fmt.Printf("addr sta.c: %p\n", sta.c)
	//
	//fmt.Println(sta.GetA())
	//fmt.Println(*sta.GetC())
	assert.NoError(t, syscall.Munmap(alloc.availableHeaps[0].heap))

}

type ttp struct {
	a [10]uint32
}

const _SizeOfTtp = unsafe.Sizeof(ttp{})

func BenchmarkDefaultAllocator_Allocate(b *testing.B) {
	b.StopTimer()

	alloc := New(create, nil)

	b.N = 1024 * 1024 * 1000
	//defer destroy()
	//defer TestReadMMapData(nil)
	//b.N = _TestMemorySize / int(_SizeOfTtp+_SizeMarker)
	b.StartTimer()
	//start := time.Now()
	for n := 0; n < b.N; n++ {
		var p = (*ttp)(alloc.Allocate(_SizeOfTtp))
		alloc.DeAllocate(unsafe.Pointer(p))
	}
	//fmt.Println(time.Now().Sub(start) / 1e3)

}

func BenchmarkGoAllocator_Allocate(b *testing.B) {
	b.N = _TestMemorySize / int(_SizeOfTtp+_SizeMarker)
	for n := 0; n < b.N; n++ {
		_ = new(ttp)
	}
}
