package synchronization

import (
	"fmt"
	"sync"
	"time"
)

/*
Work with mutex to allow synch access to the structure
*/
type SafeCounter struct {
	m   map[string]int
	mut sync.Mutex
}

func (c *SafeCounter) Inc(key string) {
	c.mut.Lock()
	c.m[key]++
	c.mut.Unlock()
}

func (c *SafeCounter) Value(key string) int {
	c.mut.Lock()
	defer c.mut.Unlock()
	return c.m[key]
}

func Experiment1() {
	counter := SafeCounter{m: make(map[string]int)}
	for i := 0; i < 1000; i++ {
		go counter.Inc("someKey")
	}

	time.Sleep(1000 * time.Millisecond)
	fmt.Println("Value =", counter.Value("someKey"))
}
