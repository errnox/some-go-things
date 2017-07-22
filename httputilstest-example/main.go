package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

func main() {
	http.HandleFunc("/", handleRoot)
	log.Fatal(http.ListenAndServe(":3044", nil))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	d, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w, fmt.Sprint(err),
			http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "%q", d)
}
