package main

import (
	"fmt"
)

func select_timer(int_chan chan int, time_out chan bool) {
	for {
		select {
		case <-int_chan:
			fmt.Println("Received msg")
		case <-time_out:
			break
		}
	}
}

func main() {
	int_chan := make(chan int)
	time_out := make(chan bool)
	go select_timer(int_chan, time_out)
	for index := 0; index < 100; index++ {
		int_chan <- index
	}
	time_out <- true
}
