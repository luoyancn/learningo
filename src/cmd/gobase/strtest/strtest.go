package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	str := "The quick brown fox jumps over the lazy dog"
	str_slice := strings.Fields(str)
	for _, val := range str_slice {
		fmt.Printf("%s\n", val)
	}

	int_str := "666lu"
	int_num, err := strconv.Atoi(int_str)
	if nil != err {
		fmt.Printf("%v\n", err)
		return
	}

	fmt.Printf("%d\n", int_num)
	new_str := strconv.Itoa(int_num)
	fmt.Printf("%s\n", new_str)
}
