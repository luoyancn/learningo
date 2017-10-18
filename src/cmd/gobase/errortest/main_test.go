package main

import (
	"fmt"
	"testing"
)

func TestBadcall(t *testing.T) {
	var err interface{} = nil
	defer func() {
		if nil == err {
			t.Errorf("Panic didn`t occured")
		}
	}()

	defer func() {
		err = recover()
		fmt.Printf("%v\n", err)
	}()
	badcall()
}
