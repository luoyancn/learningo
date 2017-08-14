package main

import (
	"fmt"
	"sync"
)

type _singleton struct {
	Name string
}

var once sync.Once
var instance *_singleton

func GetInstance() *_singleton {
	once.Do(func() {
		instance = &_singleton{
			Name: "lucifer",
		}
	})
	return instance
}

func main() {
	ins := GetInstance()
	fmt.Printf("%v\n", ins.Name)
}
