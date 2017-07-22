package main

import "fmt"

func iterator() func() int {
	a := 0
	return func() int {
		a++
		return a
	}
}

func main() {
	iter := iterator()
	fmt.Printf("Value of a: %d\n", iter())
	fmt.Printf("Value of a: %d\n", iter())
	fmt.Printf("Value of a: %d\n", iter())
}
