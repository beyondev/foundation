//+build .

package allocator

import (
	"os"
	"syscall"
	"unsafe"
)

type SegmentManager struct {
	region MappedRegion
}

func NewSegmentManager(fd int) *SegmentManager {
	syscall.Mm
}

func (s *SegmentManager) Allocate(size uintptr) unsafe.Pointer {}
func (s *SegmentManager) DeAllocate(ptr unsafe.Pointer) {}

func (s *SegmentManager) GetPointer() unsafe.Pointer {

}

type MappedRegion struct {
	base uintptr
}

func NewMappedRegion(path string, size uintptr) *MappedRegion {
	fd, err := syscall.Open(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	if _, err = syscall.Write(fd, make([]byte, size)); err != nil {
		panic(err)
	}

	data, err := syscall.Mmap(fd, 0, int(size), syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		panic(err)
	}

	if err = syscall.Close(fd); err != nil {
		panic(err)
	}

	region := &MappedRegion{}
	region.base = uintptr(unsafe.Pointer(&data[0]))

	return region
}
