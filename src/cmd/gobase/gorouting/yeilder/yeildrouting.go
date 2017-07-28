package main

import (
	"fmt"
)

type Gegerator chan interface{}

func (self *Gegerator) Next() interface{} {
	return <-(*self)
}

func Integers() Gegerator {
	yield := make(Gegerator)
	count := 0
	go func() {
		for {
			count++
			yield <- count
		}
	}()
	return yield
}

func Stringers() Gegerator {
	yield := make(Gegerator)
	str := "a"
	go func() {
		for {
			str = str + "a"
			yield <- str
		}
	}()
	return yield
}

func main() {
	resume := Integers()
	str_generator := Stringers()
	for index := 0; index < 10; index++ {
		fmt.Printf("%v\n", resume.Next())
		fmt.Printf("%v\n", str_generator.Next())
	}
}
