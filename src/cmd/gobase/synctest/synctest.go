package main

import (
	"fmt"
	"sync"
)

type syncint struct {
	mu    sync.Mutex
	count int
}

func (self *syncint) Updata() {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.count++
}

func main() {
	ch := make(chan bool, 100)
	odd_couter := &syncint{count: 0}
	oxx_couter := &syncint{count: 0}
	for i := 0; i < 100; i++ {
		go func(i int) {
			if 0 == i%2 {
				odd_couter.Updata()
				ch <- true
			} else {
				oxx_couter.Updata()
				ch <- false
			}
			<-ch
		}(i)
	}

	fmt.Println(odd_couter)
	fmt.Println(oxx_couter)
}
