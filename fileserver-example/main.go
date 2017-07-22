package main

import (
	"log"
	"net/http"
)

func main() {
	log.Fatal(http.ListenAndServe(":3045", http.FileServer(
		http.Dir("./assets"))))
}
