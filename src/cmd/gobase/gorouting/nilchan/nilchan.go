package main

import (
	"fmt"
	"math/rand"
	"time"
)

func reader(ch chan int) {
	begin := time.NewTimer(10 * time.Second)

	for {
		select {
		case i := <-ch:
			fmt.Printf("%d\n", i)
		case <-begin.C:
			ch = nil
		}
	}
}

func writer(ch chan int) {
	t := time.NewTimer(2 * time.Second)
	restart := time.NewTimer(5 * time.Second)
	savedch := ch
	for {
		select {
		case ch <- rand.Intn(42):
		case <-t.C:
			ch = nil
		case <-restart.C:
			ch = savedch
		}
	}
}

func main() {
	ch := make(chan int)
	go reader(ch)
	go writer(ch)
	time.Sleep(10 * time.Second)
}
