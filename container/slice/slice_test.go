package slice

import (
	"fmt"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	s := New(2)
	s.PushBack(1)
	s.PushBack(2)
	s.PushBack(3)
	s.PushBack(3)
	s.PushBack(3)

	fmt.Println(s, s.len, s.cap)

	s.Put(3, 4)
	s.Put(4, 5)
	fmt.Println(s, s.len, s.cap)

	//s.Put(5, 1)
}

func TestCostTime(t *testing.T) {
	const Bench = 10000000
	s := New(Bench)
	gs := make([]int, 0, Bench)

	start := time.Now()

	for n := 0; n < Bench; n++ {
		gs = append(gs, n)
	}
	fmt.Println("gosli:", time.Now().Sub(start).Nanoseconds())

	start = time.Now()
	for n := 0; n < Bench; n++ {
		s.PushBack(n)
	}
	fmt.Println("slice:", time.Now().Sub(start).Nanoseconds())

}

func BenchmarkSlice_PushBack(b *testing.B) {
	b.StopTimer()
	b.N = 10000000
	s := New(b.N)
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		s.PushBack(n)
	}

	s.Free()

}

func BenchmarkGoSlice_Append(b *testing.B) {
	b.StopTimer()
	b.N = 10000000
	s := make([]int, 0, b.N)
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		s = append(s, n)
	}
}
