package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	source, err := os.Open("copyfile.go")
	if nil != err {
		fmt.Fprintf(os.Stderr,
			"An error occured when opening file: %v\n", err)
		os.Exit(2)
	}
	defer source.Close()

	dest, err := os.OpenFile("dest.txt", os.O_CREATE|os.O_WRONLY, 0644)
	if nil != err {
		fmt.Fprintf(os.Stderr,
			"An error occured when creating file: %v\n", err)
		os.Exit(2)
	}

	defer dest.Close()
	_, err = io.Copy(dest, source)
	if nil != err {
		fmt.Fprintf(os.Stderr,
			"An error occured when copying file: %v\n", err)
		os.Exit(2)
	}
}
