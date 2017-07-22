package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
)

var port string
var defined bool

func init() {
	port, defined = os.LookupEnv("PORT")
	if !defined {
		port = "4321"
	}
}

func main() {
	http.HandleFunc("/", handleHTTP)

	http.HandleFunc("/info", handleInfo)
	http.HandleFunc("/info/", handleInfo)

	http.HandleFunc("/about", handleAbout)
	http.HandleFunc("/about/", handleAbout)

	http.HandleFunc("/agent", handleAgent)
	http.HandleFunc("/agent/", handleAgent)

	http.HandleFunc("/redirect", handleRedirect)
	http.HandleFunc("/redirect/", handleRedirect)

	http.HandleFunc("/file", serveMain)
	http.HandleFunc("/file/", serveMain)

	fmt.Printf("Listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func listenAndServe() {
	s := http.Server{
		Handler: http.HandlerFunc(handleHTTP),
		Addr:    ":3044",
	}
	log.Fatal(s.ListenAndServe())
}

func handleHTTP(w http.ResponseWriter, r *http.Request) {
	s := ""
	for i := 0; i < 10; i++ {
		s = s + fmt.Sprintf("%d: %s\n", i, html.EscapeString(
			r.URL.Path))
	}
	fmt.Fprintf(w, "%s", s)
}

func handleInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Info", "This is some info.")
	fmt.Fprintf(w, "Here is some information.\n")
}

func handleAbout(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the About page.\n")
}

func handleAgent(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s\n", r.UserAgent())
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/redirected", http.StatusMovedPermanently)
}

func serveMain(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "main.go")
}
