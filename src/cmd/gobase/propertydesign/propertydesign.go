package main

import (
	"designpatterns"
	"fmt"

	uuid "github.com/satori/go.uuid"
)

func main() {
	config := designpatterns.NewConfig("guest", "/home/guest")

	root_config := config.WithUser("root").WithWorkDir("/root")

	fmt.Printf("guest config %v\n", config)
	fmt.Printf("root config %v\n", root_config)

	fmt.Printf("Uuid is %v\n", uuid.NewV4().String())
}
