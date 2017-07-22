package main

import (
	"flag"
	"fmt"
)

var repetitions int
var helpmsg bool

func init() {
	const (
		repetitionsDefault = 1
		repetitionsHelp    = "Number of repetitions to print"

		helpmsgDefault = false
		helpmsgHelp    = "Print a help message"
	)

	flag.IntVar(&repetitions, "repeat", repetitionsDefault,
		repetitionsHelp)
	flag.IntVar(&repetitions, "r", repetitionsDefault,
		repetitionsHelp)

	flag.BoolVar(&helpmsg, "help", helpmsgDefault, helpmsgHelp)
	flag.BoolVar(&helpmsg, "h", helpmsgDefault, helpmsgHelp)
}

func main() {
	flag.Parse()
	if helpmsg == true || flag.NArg() == 0 {
		flag.PrintDefaults()
	} else {
		printRepetitions()
	}

}

func printRepetitions() {
	for i := 0; i < repetitions; i++ {
		fmt.Printf("%v\n", flag.Arg(flag.NArg()-1))
	}
}
