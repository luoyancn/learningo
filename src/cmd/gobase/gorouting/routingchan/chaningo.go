package main

import (
	"fmt"
)

func f(left chan int, right chan int) {
	left <- 1 + <-right
}

func main() {
	leftmost := make(chan int)
	var left, right chan int = nil, leftmost
	for i := 0; i < 1000; i++ {
		left, right = right, make(chan int)
		go f(left, right)
	}

	right <- 0
	x := <-leftmost
	fmt.Println(x)
}
