package main

import (
	"designpatterns"
	"fmt"
	"os"
)

func main() {
	sender := &designpatterns.Sender{}

	jsonmsg, err := sender.BuildMessage(&designpatterns.JSONMessageBuilder{})
	if nil != err {
		fmt.Printf("ERROR:%v\n", err)
		os.Exit(2)
	}
	fmt.Println(string(jsonmsg.Body))

	xmlmsg, err := sender.BuildMessage(&designpatterns.XMLMessageBuilder{})
	if nil != err {
		fmt.Printf("ERROR:%v\n", err)
		os.Exit(2)
	}
	fmt.Println(string(xmlmsg.Body))
}
