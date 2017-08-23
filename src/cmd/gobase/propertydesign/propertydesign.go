package main

import (
	"designpatterns"
	"fmt"
)

func main() {
	config := designpatterns.NewConfig("guest", "/home/guest")

	root_config := config.WithUser("root").WithWorkDir("/root")

	fmt.Printf("guest config %v\n", config)
	fmt.Printf("root config %v\n", root_config)
}
