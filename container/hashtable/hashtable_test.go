package hashtable

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	h := New(0)
	defer h.Free()
	h.Put(1, 1) //1
	h.Put(2, 1) //2
	h.Put(2, 2)
	h.Put(3, 2)  //4
	h.Put(4, 2)  //4
	h.Put(5, 2)  //8
	h.Put(6, 2)  //8
	h.Put(61, 2) //8
	h.Put(62, 2) //8
	assert.Equal(t, 8, h.Cap())
	h.Put(63, 2) //16
	assert.Equal(t, 16, h.Cap())
	h.Put(163, 2)
	h.Put(263, 2)
	h.Put(64, 2)
	h.Put(32, 2)

	m := map[V]V{
		1:   1,
		2:   2,
		3:   2,
		4:   2,
		5:   2,
		6:   2,
		32:  2,
		61:  2,
		62:  2,
		63:  2,
		64:  2,
		163: 2,
		263: 2,
	}

	h.Each(func(key K, value V) {
		assert.Equal(t, value, m[key])
	})

	fmt.Println(m)
	fmt.Println(h)

	fmt.Println("cap:", h.Cap())
	h.Remove(61)
	h.Remove(62)
	h.Remove(63)
	e := (*Buckets)(h.b.Get()).At(0)
	h.Remove(64)
	e = (*Buckets)(h.b.Get()).At(0)
	fmt.Println(e)
	h.Remove(32)

	fmt.Println(h)
}

func TestNewBuckets(t *testing.T) {
	b := NewBuckets(0)
	b.Free()
}

func TestBuckets_At(t *testing.T) {
	b := NewBuckets(2)
	b.Put(1, &Entry{key: 1, value: 2})
	//e := b.At(1)
	//e.key = 1
	//e.value = 2
	ee := b.At(1)
	fmt.Println(*ee)
	b.Free()
}

func TestIterator_Next(t *testing.T) {
	h := New(3)
	h.Put(1, 1)
	h.Put(2, 1)
	h.Put(2, 2)
	h.Put(3, 2)

	for it := h.Begin(); !it.IsEnd(); it.Next() {
		fmt.Print("(", it.Key(), ":", it.Value(), "),")
	}
	fmt.Println()

	h.Free()
}

func Test_CostTime(t *testing.T) {
	const Bench = 1000000
	h := New(0)
	m := make(map[int]int, 0)

	start := time.Now()
	for n := 0; n < Bench; n++ {
		h.Put(n, n)
	}
	fmt.Println("table-put", time.Now().Sub(start).Nanoseconds())

	start = time.Now()
	for n := 0; n < Bench; n++ {
		m[n] = n
	}
	fmt.Println("gomap-put", time.Now().Sub(start).Nanoseconds())

	start = time.Now()
	for n := h.Begin(); !n.IsEnd(); n.Next() {
	}
	fmt.Println("table-loop", time.Now().Sub(start).Nanoseconds())

	start = time.Now()
	for range m {
	}
	fmt.Println("gomap-loop", time.Now().Sub(start).Nanoseconds())

	start = time.Now()
	for n := 0; n < Bench; n++ {
		h.Remove(n)
	}
	fmt.Println("table-rm", time.Now().Sub(start).Nanoseconds())

	start = time.Now()
	for n := 0; n < Bench; n++ {
		delete(m, n)
	}
	runtime.GC()
	fmt.Println("gomap-rm", time.Now().Sub(start).Nanoseconds())

	h.Free()

}

func BenchmarkBuckets_Put(b *testing.B) {
	b.StopTimer()
	b.N = 1000000
	h := New(b.N)
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		h.Put(n, n)
	}
	h.Free()
}

func BenchmarkMap_Put(b *testing.B) {
	b.StopTimer()
	b.N = 1000000
	h := make(map[int]int, b.N)
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		h[n] = n
	}
}
