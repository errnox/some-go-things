package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
)

const (
	RELOAD = false
)

func main() {
	printManpages()
	fmt.Println("It works.")
}

func printManpages() {
	manpage := FSMustString(RELOAD, "/assets/manpage.txt")
	fmt.Println(manpage)

	bashpage := FSMustString(RELOAD, "/assets/bashpage.txt")
	fmt.Println(bashpage)

	gitpage := FSMustString(RELOAD, "/assets/gitpage.txt")
	fmt.Println(gitpage)
}

func executeCommand() {
	_, err := exec.Command("feh",
		"/usr/share/images/fluxbox/fluxbox.png").Output()
	if err != nil {
		log.Fatal(err)
	}
}

func writeImageFile() {
	img := FSMustByte(RELOAD, "/assets/image.png")
	err := ioutil.WriteFile("output-image.png", img, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
