package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	_       = iota
	BIT int = iota * 2
	HALFBYTE
	BYTE
	DOUBLEBYTE
)

type car struct {
	topSpeed float32
	numTires int
	size     string
}

func init() {
	// fmt.Println("Initializing...")
	// fmt.Println("Done initializing")
}

func main() {
	printDuration()
}

func printPrintf() {
	x := 3.3
	fmt.Printf("%v\n", int64(x))
}
func printMap() {
	colors := map[string]uint32{
		"red":   0xff0000,
		"green": 0x00ff00,
		"blue":  0x0000ff,
		"white": 0xffffff,
		"black": 0x000000,
	}
	fmt.Printf("%v\n", colors)
	fmt.Printf("red:   %x\n", colors["red"])
	fmt.Printf("green: %x\n", colors["green"])
}

func (t car) String() string {
	return fmt.Sprintf("Car\n---\nspeed: %.2f\ntires: %d\nsize:  %s",
		t.topSpeed, t.numTires, t.size)
}

func printCustomType() {
	redCar := car{topSpeed: 180.0, numTires: 4, size: "medium"}

	fmt.Printf("%v\n", redCar)

	a := []int{111, 222, 333, 444}
	a = append(a, []int{555, 666, 777, 888, 999}...)
	fmt.Println(a)
}

func printIotaConstants() {
	fmt.Println("       BIT: ", BIT)
	fmt.Println("  HALFBYTE: ", HALFBYTE)
	fmt.Println("      BYTE: ", BYTE)
	fmt.Println("DOUBLEBYTE: ", DOUBLEBYTE)
}

func printInterface() {
	cars := []car{}
	for i := 0; i < 15; i++ {
		cars = append(cars,
			car{topSpeed: 100.0 + float32(i), numTires: i})
	}
	fmt.Printf("%v\n", cars)
	fmt.Printf("\nThere are %d cars.\n", len(cars))
}

func printIoutilReadAll() {
	r := strings.NewReader("one\ntwo\nthree\nfour\n")
	x, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", x)
}

func printIoutilReadDir() {
	d, err := ioutil.ReadDir(filepath.Join(os.Getenv("HOME"),
		".config"))
	if err != nil {
		log.Fatal(err)
	}
	s := ""
	for _, f := range d {
		if f.IsDir() {
			s = "/"
		} else {
			s = ""
		}
		fmt.Printf("%10d  %s%s\n", f.Size(), f.Name(), s)
	}
}

func printIoutilPrintFile() {
	f, err := ioutil.ReadFile("./main.go")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", f)
}

func printDuration() {
	s := "02h50m03s"
	d, err := time.ParseDuration(s)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n\n", s)
	fmt.Printf("Hours:       %v\n", d.Hours())
	fmt.Printf("Minutes:     %v\n", d.Minutes())
	fmt.Printf("Seconds:     %v\n", d.Seconds())
	fmt.Printf("Nanoseconds: %v\n", d.Nanoseconds())
	fmt.Printf("\nString:      %v\n", d.String())
}
