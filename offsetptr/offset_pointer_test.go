package offsetptr

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
	"unsafe"
)

func TestPointer_Get(t *testing.T) {
	p := Pointer{}
	pint := new(int)
	*pint = 123
	p.Set(unsafe.Pointer(pint))

	assert.Equal(t, 123, *(*int)(p.Get()))

	pp := Pointer{}
	pp.Set(unsafe.Pointer(pint))

	assert.Equal(t, 123, *(*int)(pp.Get()))

	p1 := Pointer{} //<Pointer>
	pr := new(int)
	*pr = 3

	p2 := NewPointer(unsafe.Pointer(pr))
	p1.Set(unsafe.Pointer(p2))

	fmt.Println(*(*int)((*Pointer)(p1.Get()).Get()))

}

func Test_Offset(t *testing.T) {
	a := uintptr(0x100)
	b := uintptr(0x101)
	off := a - b
	c := b + off
	fmt.Println(off, c)
}

func Test_OffsetToRaw(t *testing.T) {
	type raw struct {
		a    int
		next Pointer
	}

	raw1 := raw{a: 10}
	raw1.next.Set(nil)

	raw2 := raw{a: 20}
	raw2.next.Set(unsafe.Pointer(&raw1))

	off := Pointer{}
	off.Set(unsafe.Pointer(&raw2))

	r2 := (*raw)(off.Get())
	r1 := (*raw)(r2.next.Get())

	fmt.Println(r1.a)
	fmt.Println(r2.a)
}

func Test_address(t *testing.T) {
	type N struct {
		a int
	}

	//n := N{}
	//fmt.Printf("%p\n", &n)
	//fmt.Println(&n.a)

	p1 := Pointer{}
	p2 := Pointer{}

	//fmt.Printf("%p\n", &p2)
	p2.Get()
	p1.Forward(&p2)

	p := new(N)
	fmt.Println(uintptr(unsafe.Pointer(p)))
	fmt.Println(uintptr(unsafe.Pointer(&p)))
	fmt.Println(*(*uintptr)(unsafe.Pointer(p)))
}

func Test_refrence(t *testing.T) {
	type sd struct {
		p Pointer
	}

	pp := Pointer{}
	pp.Set(unsafe.Pointer(new(int)))
	*(*int)(pp.Get()) = 100

	s := sd{}
	s.p.Forward(&pp)

	fmt.Println(*(*int)(s.p.Get()))
}

type offT struct {
	x *int
	y *yT
}

func NewOffT() *offT {
	t := new(offT)
	t.x = new(int)
	t.y = NewYT()
	*t.x = 1
	t.y.a = 2
	runtime.SetFinalizer(t, (*offT).Free)
	return t
}

func (t *offT) Free() {
	fmt.Println("offT deleted!")
}

type yT struct {
	a int
}

func NewYT() *yT {
	y := new(yT)
	runtime.SetFinalizer(y, (*yT).Free)
	return y
}

func (*yT) Free() {
	fmt.Println("yT freed")
}

var persist = struct {
	p    *offT
	offP Pointer
}{}

func callGC() {
	d := NewOffT()
	persist.p = d
	//runtime.KeepAlive(&d)
	//persist.offP.Set(unsafe.Pointer(d))
	runtime.GC()
}

func Test_gc(t *testing.T) {
	callGC()
	runtime.GC()

	//for n := 0; n < 10000000; n++ {
	//	_ = NewOffT()
	//}

	fmt.Println(*persist.p.x)
	//fmt.Println(*(*offT)(persist.offP.Get()).x)
	//(*offT)(persist.offP.Get()).free()

}
