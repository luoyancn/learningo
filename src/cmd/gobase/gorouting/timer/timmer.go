package main

import (
	"fmt"
	"time"
)

func main() {
	tick := time.Tick(time.Millisecond * 10)
	boom := time.After(time.Millisecond * 100)
	for {
		select {
		case <-tick:
			fmt.Println("tick.")
		case <-boom:
			fmt.Println("BOOM!")
			return
		default:
			fmt.Println("       .")
			time.Sleep(time.Millisecond * 10)
		}
	}
}
