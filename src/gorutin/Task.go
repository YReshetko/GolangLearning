package gorutin

import (
	"fmt"

	"golang.org/x/tour/tree"
)

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	if t != nil {
		Walk(t.Left, ch)
		ch <- t.Value
		Walk(t.Right, ch)
	}
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	c1 := make(chan int)
	c2 := make(chan int)
	go Walk(t1, c1)
	go Walk(t2, c2)
	result := true
	for i := 0; i < 10; i++ {
		if <-c1 != <-c2 {
			result = false
		}
	}
	return result
}

func Task() {

	// Test walk
	ch := make(chan int)
	go Walk(tree.New(2), ch)
	for i := 0; i < 10; i++ {
		fmt.Println(<-ch)
	}

	// Test same
	fmt.Printf("Test equal trees: %v\n", Same(tree.New(1), tree.New(1)))
	fmt.Printf("Test different trees: %v\n", Same(tree.New(1), tree.New(2)))
}
