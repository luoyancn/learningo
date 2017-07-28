package main

import (
	"fmt"
	"runtime"
)

func teller(deposits chan int, balances chan int) {
	var balance int
	for {
		select {
		case amount := <-deposits:
			balance += amount
		case balances <- balance:
		}
	}
}

func Deposit(amount int, deposits chan int) {
	deposits <- amount
}

func Balance(balances chan int) int {
	return <-balances
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	deposits := make(chan int)
	balances := make(chan int)
	go teller(deposits, balances)

	done := make(chan struct{})
	go func() {
		Deposit(200, deposits)
		done <- struct{}{}
	}()

	go func() {
		Deposit(200, deposits)
		done <- struct{}{}
	}()

	<-done
	<-done
	fmt.Println(Balance(balances))
}
