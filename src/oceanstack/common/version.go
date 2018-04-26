package common

import (
	"fmt"
)

var Version = "No Version Provided"
var Buildstamp = "No build stamp Provided"
var Githash = "No githash Provided"

func Versions(product string) {
	fmt.Printf("%s Version: %s\n", product, Version)
	fmt.Printf("%s buildstamp: %s\n", product, Buildstamp)
	fmt.Printf("%s githash: %s\n", product, Githash)
}
