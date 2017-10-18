package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	var input_reader *bufio.Reader
	input_reader = bufio.NewReader(os.Stdin)
	fmt.Println("Please enter some input:")
	input, err := input_reader.ReadString('\n')
	if nil != err {
		fmt.Fprintf(os.Stderr, "Fail to read from stdin: %v\n", err)
		os.Exit(2)
	}

	fmt.Printf("The input was :%s\n", input)
	input_file, err := os.Open("bufiotest.go")
	if nil != err {
		fmt.Fprintf(os.Stderr,
			"An error occured when openging the input file: %v\n", err)
		os.Exit(2)
	}

	defer input_file.Close()
	input_reader = bufio.NewReader(input_file)

	output_file, err := os.OpenFile("bufio.txt", os.O_CREATE|os.O_WRONLY, 0644)
	if nil != err {
		fmt.Fprintf(os.Stderr,
			"An error occured when creating file: %v\n", err)
		os.Exit(2)
	}

	defer output_file.Close()
	output_writer := bufio.NewWriter(output_file)
	for {
		input_str, err := input_reader.ReadString('\n')
		if io.EOF == err {
			output_writer.Flush()
			return
		}
		//output_file.WriteString(input_str)
		output_writer.WriteString(input_str)
	}

}
