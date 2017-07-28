package main

import (
	"fmt"
)

type Any interface{}
type EvalFunc func(Any) (Any, Any)

func BuildLazyEvaluator(evenFunc EvalFunc, initState Any) func() Any {
	retValChan := make(chan Any)
	loopFunc := func() {
		var retVal Any
		var actState Any = initState
		for {
			retVal, actState = evenFunc(actState)
			retValChan <- retVal
		}
	}

	retFunc := func() Any {
		return <-retValChan
	}

	go loopFunc()
	return retFunc
}

func BuildIntLazyEvaluator(evalFunc EvalFunc, initState Any) func() int {
	ef := BuildLazyEvaluator(evalFunc, initState)
	return func() int {
		return ef().(int)
	}
}

func main() {
	evenFunc := func(state Any) (Any, Any) {
		os := state.(int)
		ns := os + 2
		return os, ns
	}

	even := BuildIntLazyEvaluator(evenFunc, 0)
	for i := 0; i < 10; i++ {
		fmt.Printf("%vth even:%v\n", i, even())
	}
}
