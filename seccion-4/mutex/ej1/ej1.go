package main

import (
	"fmt"
	"sync"
)

var x = 0
var m sync.Mutex

func increment(wg *sync.WaitGroup) {
	for i := 0; i < 100; i++ {
		m.Lock()
		x = x + 1
		m.Unlock()
	}
	defer wg.Done()
}
func main() {
	var w sync.WaitGroup
	w.Add(2)
	go increment(&w)
	go increment(&w)
	w.Wait()
	fmt.Println("final value of x", x)
}
