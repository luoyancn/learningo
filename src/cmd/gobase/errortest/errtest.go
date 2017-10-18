package main

import (
	"fmt"
)

func badcall() {
	panic("bad end")
}

func test() {
	defer func() {
		if e := recover(); nil != e {
			fmt.Printf("Panic :%v\n", e)
		}
	}()
	badcall()
	fmt.Printf("After bad call\n")
}

func main() {
	fmt.Printf("Calling test \n")
	test()
	fmt.Printf("Test Completed\n")
}
