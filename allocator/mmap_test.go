package allocator

import (
	"fmt"
	"github.com/gotests/got/shm/entity"
	"log"
	"os"
	"syscall"
	"testing"
	"unsafe"
)

func TestWrite(t *testing.T) {

	f, err := syscall.Open(_TestMemoryFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if nil != err {
		log.Fatalln(err)
	}

	if _, err := syscall.Write(f, make([]byte, _TestMemorySize)); nil != err {
		log.Fatalln(err)
	}
	// extend file
	//if _, err := f.WriteAt([]byte{byte(0)}, _TestMemorySize); nil != err {
	//	log.Fatalln(err)
	//}

	data, err := syscall.Mmap(f, 0, _TestMemorySize, syscall.PROT_WRITE, syscall.MAP_SHARED|syscall.MAP_COPY)
	//data, err := syscall.Mmap(int(f.Fd()), 0, _TestMemorySize, syscall.PROT_WRITE, syscall.MAP_SHARED|syscall.MAP_COPY)
	//data, err := unix.Mmap(int(f.Fd()), 0, 1<<8, syscall.PROT_WRITE, syscall.MAP_SHARED|syscall.MAP_COPY)
	if nil != err {
		log.Fatalln(err)
	}

	defer func() {
		if err := syscall.Munmap(data); nil != err {
			log.Fatalln("@@", err)
		}
	}()

	if err := syscall.Close(f); nil != err {
		log.Fatalln(err)
	}

	qq := (*entity.Share)(unsafe.Pointer(&data[0]))

	qq.Id = 100

	//qq.Idp = (*uint32)(unsafe.Pointer(&data[entity.SizeOfShare]))
	//
	//*qq.Idp = 16

	//fmt.Printf("%v\n", data)
	//fmt.Printf("qq: %p\n", qq)
	//fmt.Printf("qq.id: %p\n", qq.Idp)
}

func TestRead(t *testing.T) {

	f, err := syscall.Open(_TestMemoryFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if nil != err {
		log.Fatalln(err)
	}

	data, err := syscall.Mmap(f, 0, _TestMemorySize, syscall.PROT_WRITE, syscall.MAP_SHARED)
	if nil != err {
		log.Fatalln(err)
	}

	defer func() {
		if err := syscall.Munmap(data); nil != err {
			log.Fatalln("@@", err)
		}
	}()

	if err := syscall.Close(f); nil != err {
		log.Fatalln(err)
	}

	//
	qq := (*entity.Share)(unsafe.Pointer(&data[0]))

	fmt.Printf("data: %p\n", data)
	fmt.Printf("qq: %p\n", qq)
	//fmt.Printf("qq.id: %p\n", qq.Idp)

	fmt.Println(*qq)
	//fmt.Println(*qq.Idp)

	//idps := (*int)(unsafe.Pointer(&data[entity.SizeOfShare]))
	//fmt.Println(*idps)
	//fmt.Printf("%p\n", idps)

}

type shared struct {
	id  uint64
	id2 uint64
	idp uintptr `*uint64`
}

func TestWrite2(t *testing.T) {

	f, err := syscall.Open(_TestMemoryFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if nil != err {
		log.Fatalln(err)
	}

	if _, err := syscall.Write(f, make([]byte, _TestMemorySize)); nil != err {
		log.Fatalln(err)
	}
	// extend file
	//if _, err := f.WriteAt([]byte{byte(0)}, _TestMemorySize); nil != err {
	//	log.Fatalln(err)
	//}

	data, err := syscall.Mmap(f, 0, _TestMemorySize, syscall.PROT_WRITE, syscall.MAP_SHARED|syscall.MAP_COPY)
	//data, err := syscall.Mmap(int(f.Fd()), 0, _TestMemorySize, syscall.PROT_WRITE, syscall.MAP_SHARED|syscall.MAP_COPY)
	//data, err := unix.Mmap(int(f.Fd()), 0, 1<<8, syscall.PROT_WRITE, syscall.MAP_SHARED|syscall.MAP_COPY)
	if nil != err {
		log.Fatalln(err)
	}

	defer func() {
		if err := syscall.Munmap(data); nil != err {
			log.Fatalln("@@", err)
		}
	}()

	if err := syscall.Close(f); nil != err {
		log.Fatalln(err)
	}

	base := uintptr(unsafe.Pointer(&data[0]))

	//*(*uint64)(unsafe.Pointer(base + 0)) = 100
	//*(*uint64)(unsafe.Pointer(base + unsafe.Sizeof(uint64(0)))) = 120
	//*(*uintptr)(unsafe.Pointer(base + 2*unsafe.Sizeof(uint64(0)))) = 2*unsafe.Sizeof(uint64(0)) + unsafe.Sizeof(uintptr(0))
	//
	//*(*uint64)(unsafe.Pointer(base + 2*unsafe.Sizeof(uint64(0)) + unsafe.Sizeof(uintptr(0)))) = 16

	s := (*shared)(unsafe.Pointer(base + 0))
	s.id = 100
	s.id2 = 120
	s.idp = unsafe.Sizeof(shared{})

	idp := (*uint64)(unsafe.Pointer(base + s.idp))
	*idp = 16

	fmt.Println(s)
}

func TestRead2(t *testing.T) {

	f, err := syscall.Open(_TestMemoryFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if nil != err {
		log.Fatalln(err)
	}

	data, err := syscall.Mmap(f, 0, _TestMemorySize, syscall.PROT_WRITE, syscall.MAP_SHARED)
	if nil != err {
		log.Fatalln(err)
	}

	defer func() {
		if err := syscall.Munmap(data); nil != err {
			log.Fatalln("@@", err)
		}
	}()

	if err := syscall.Close(f); nil != err {
		log.Fatalln(err)
	}

	base := uintptr(unsafe.Pointer(&data[0]))
	shared := (*shared)(unsafe.Pointer(base + 0))
	fmt.Println(shared)

	id := shared.id
	id2 := shared.id2
	idp := (*uint64)(unsafe.Pointer(base + shared.idp))

	fmt.Println(id, id2, *idp)
}

type sofs struct {
	a uint64
	b [500]uint64
	c []byte
	d string
	e map[int]string
}

func BenchmarkSizeof(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = unsafe.Sizeof(sofs{})
	}
}

const sof = unsafe.Sizeof(sofs{})

func BenchmarkConstSizeof(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = sof
	}
}
