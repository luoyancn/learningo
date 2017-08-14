package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	now := time.Now()
	fmt.Println(now.Format(time.RFC1123))
	fmt.Println(now.Format(time.RFC1123Z))
	fmt.Println(now.Format(time.RFC3339))
	fmt.Println(now.Format(time.RFC3339Nano))
	fmt.Println(now.Format(time.RFC822))
	fmt.Println(now.Format(time.RFC822Z))
	fmt.Println(now.Format(time.RFC850))

	begin, err := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
	if nil != err {
		fmt.Println(err)
		os.Exit(2)
	}
	fmt.Println(begin)
}
