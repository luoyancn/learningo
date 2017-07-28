package main

import (
	"fmt"
)

func generator1(ch chan int) {
	for i := 2; ; i++ {
		ch <- i
	}
}

func filter1(in, out chan int, prime int) {
	for {
		i := <-in
		if 0 != i%prime {
			out <- i
		}
	}
}

func sieve1() {
	begin := make(chan int)
	go generator1(begin)
	for i := 0; i < 100; i++ {
		prime := <-begin
		fmt.Printf("Prime is %d\n", prime)
		mid_chan := make(chan int)
		go filter1(begin, mid_chan, prime)
		begin = mid_chan
	}
}

func generator2() chan int {
	ch := make(chan int)
	go func() {
		for i := 2; ; i++ {
			ch <- i
		}
	}()
	return ch
}

func filter2(in chan int, prime int) chan int {
	out := make(chan int)
	go func() {
		for {
			if i := <-in; 0 != i%prime {
				out <- i
			}
		}
	}()

	return out
}

func sieve2() chan int {
	out := make(chan int)
	go func() {
		ch := generator2()
		for {
			prime := <-ch
			ch = filter2(ch, prime)
			out <- prime
		}
	}()

	return out
}

func emit(ch chan string) {
	words := []string{"a", "b", "c"}
	for _, word := range words {
		ch <- word
	}
	close(ch)
}
func main() {
	sieve1()

	primes := sieve2()
	for i := 0; i < 10; i++ {
		fmt.Printf("The second is %d\n", <-primes)
	}

	wordchan := make(chan string)
	go emit(wordchan)
	for word := range wordchan {
		fmt.Printf("%s\n", word)
	}
}
