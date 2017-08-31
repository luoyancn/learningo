package main

import (
	"fmt"
	"time"
)

func emit(chan_chan chan chan string, done chan bool) {
	word_chan := make(chan string)
	chan_chan <- word_chan
	defer close(word_chan)

	words := []string{"zhangjl", "luoyan", "zhangzzh"}
	time_out := time.NewTimer(1 * time.Second)
	i := 0
	for {
		select {
		case word_chan <- words[i]:
			i += 1
			if len(words) == i {
				i = 0
			}
		case <-done:
			done <- true
			return
		case <-time_out.C:
			return
		}
	}
}

func main() {
	chan_chan := make(chan chan string)
	done := make(chan bool)

	go emit(chan_chan, done)

	word_chan := <-chan_chan
	for word := range word_chan {
		fmt.Printf("%s", word)
	}
}
