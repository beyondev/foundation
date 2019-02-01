package redblacktree

import (
	"fmt"
	"testing"
)

func TestTree_Insert(t *testing.T) {
	tree := NewWithIntComparator(false)
	tree.Insert(1,2)
	fmt.Println(tree)
	tree.Insert(2,3)
	fmt.Println(tree)
	tree.Insert(3,5)
	fmt.Println(tree)
}