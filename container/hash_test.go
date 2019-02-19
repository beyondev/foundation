package container

import (
	"fmt"
	"testing"
)

func Test_hash(t *testing.T) {
	mp := map[*int]int{}
	mp[nil] = 1

}

func TestHash(t *testing.T) {
	fmt.Println(Hash(15))
	fmt.Println(Hash(15))

	fmt.Println(Hash("abc"))
	fmt.Println(Hash("abc"))

	type st struct {
		a uint32
		b bool
		c *int
	}

	intp := new(int)

	fmt.Println(Hash(st{1023, true, intp}))
	fmt.Println(Hash(st{1023, true, intp}))

	type sts struct {
		a uint32
		c string
	}

	s1 := "abc"
	s2 := "abc"

	fmt.Println(Hash(sts{1023, s1}))
	fmt.Println(Hash(sts{1023, s2}))

	fmt.Println(Hash([4]uint64{123, 21, 36, 49}))
	fmt.Println(Hash([4]uint64{123, 21, 36, 49}))

	fmt.Println(Hash(nil))
	fmt.Println(Hash((*int)(nil)))
	fmt.Println(Hash((*uint32)(nil)))
}
