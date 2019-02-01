package empty_index

import (
	"github.com/eosspark/eos-go/common/container"
	"github.com/eosspark/eos-go/common/container/allocator"
	"log"
	"os"
	"syscall"
)

//go:generate go install "github.com/eosspark/eos-go/common/container/multiindex/multi_index_container/..."
//go:generate go install "github.com/eosspark/eos-go/common/container/multiindex/ordered_index/..."

const (
	_TestMemoryFilePath = "/tmp/data/mmap.bin"
	_TestMemorySize     = 1024 * 1024 * 1024 * 4
)

var create = func() []byte {
	f, err := os.OpenFile(_TestMemoryFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if nil != err {
		log.Fatalln(err)
	}

	// extend file
	if _, err := f.WriteAt([]byte{byte(0)}, int64(_TestMemorySize)); nil != err {
		log.Fatalln(err)
	}

	data, err := syscall.Mmap(int(f.Fd()), 0, _TestMemorySize, syscall.PROT_WRITE, syscall.MAP_SHARED)
	if nil != err {
		log.Fatalln(err)
	}

	if err := f.Close(); nil != err {
		log.Fatalln(err)
	}

	return data
}

var alloc = allocator.NoAlloc

type item struct {
	id  uint32
	num uint64
}

//go:generate gotemplate "github.com/eosspark/eos-go/common/container/multiindex/multi_index_container" EpIndex(ByID,ByIDNode,item,alloc)
func (m *EpIndex) GetByID() *ByID {
	return (*ByID)(m.super.Get())
}

//go:generate gotemplate "github.com/eosspark/eos-go/common/container/multiindex/ordered_index" ByID(EpIndex,EpIndexNode,EpIndexBase,EpIndexBaseNode,item,uint32,ByIdKeyFunc,ByIdCompare,false,alloc)
var ByIdKeyFunc = func(n item) uint32 { return n.id }
var ByIdCompare = container.UInt32Comparator

//go:generate go build .
