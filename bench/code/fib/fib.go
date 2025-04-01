package main

import (
	"C" // without this line, Go toolchain does not work. It is a bug in the Golang compiler
	"fmt"
)

func fib(n uint) uint {
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return fib(n-1) + fib(n-2)
	}
}

func main() {
	n := uint(40)
	fmt.Println(fib(n))
}
