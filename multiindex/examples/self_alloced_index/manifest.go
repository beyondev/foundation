package self_alloced_index

import (
	"github.com/eosspark/eos-go/common/container"
	"github.com/eosspark/eos-go/common/container/allocator"
	"log"
	"os"
	"syscall"
	"unsafe"
)

const _TestMemorySize = 1024 * 1024 * 1024 * 4
const _TestMemoryFilePath = "/tmp/data/mmap.bin"

var fid int
var defaultAlloc = allocator.NewDefaultAllocator(func() []byte {
	f, err := os.OpenFile(_TestMemoryFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if nil != err {
		log.Fatalln(err)
	}

	// extend file
	if _, err := f.WriteAt([]byte{byte(0)}, int64(_TestMemorySize)); nil != err {
		log.Fatalln(err)
	}

	fid = int(f.Fd())

	data, err := syscall.Mmap(fid, 0, _TestMemorySize, syscall.PROT_WRITE, syscall.MAP_SHARED|syscall.MAP_COPY)
	//data, err := unix.Mmap(int(f.Fd()), 0, 1<<8, syscall.PROT_WRITE, syscall.MAP_SHARED|syscall.MAP_COPY)
	if nil != err {
		log.Fatalln(err)
	}

	if err := f.Close(); nil != err {
		log.Fatalln(err)
	}

	return data
}, nil)

type SharedClass struct {
	id  int
	num uint32
}

const _SizeOfSharedClass = unsafe.Sizeof(SharedClass{})

func NewSharedClass() *SharedClass {
	return (*SharedClass)(defaultAlloc.Allocate(_SizeOfSharedClass))
}

type ValueType = SharedClass

//go:generate go install "github.com/eosspark/eos-go/common/container/"
//go:generate go install "github.com/eosspark/eos-go/common/container/allocator/"
//go:generate go install "github.com/eosspark/eos-go/common/container/multiindex/"
//go:generate go install "github.com/eosspark/eos-go/common/container/multiindex/ordered_index/..."
//go:generate go install "github.com/eosspark/eos-go/common/container/multiindex/multi_index_container/..."

//go:generate gotemplate "github.com/eosspark/eos-go/common/container/multiindex/multi_index_container" TestIndex(ById,ByIdNode,ValueType,defaultAlloc)
// go:generate gotemplate "github.com/eosspark/eos-go/common/container/multiindex/multi_index_container" TestIndex(ById,ByIdNode,ValueType,allocator.NoAlloc)

//go:generate gotemplate "github.com/eosspark/eos-go/common/container/multiindex/ordered_index" ById(TestIndex,TestIndexNode,TestIndexBase,TestIndexBaseNode,ValueType,int,ByIdKeyFunc,ByIdCompare,true,defaultAlloc)
// go:generate gotemplate "github.com/eosspark/eos-go/common/container/multiindex/ordered_index" ById(TestIndex,TestIndexNode,TestIndexBase,TestIndexBaseNode,ValueType,int,ByIdKeyFunc,ByIdCompare,true,allocator.NoAlloc)
var ByIdKeyFunc = func(n SharedClass) int { return n.id }
var ByIdCompare = container.IntComparator
//go:generate go build .
///go:generate gotemplate "github.com/eosspark/eos-go/common/container/multiindex/ordered_index" ByPrev(TestIndex,TestIndexNode,ByNum,ByNumNode,ValueType,int,ByPrevKeyFunc,ByPrevCompare,true)
