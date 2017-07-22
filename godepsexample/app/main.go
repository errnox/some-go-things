package main

import (
	"fmt"
	maths "godepsexample/app/math"
	"godepsexample/app/strings"
)

func main() {
	fmt.Println("Hello there!")
	strings.Hline()
	fmt.Println(maths.Sum(3, 4))
}
