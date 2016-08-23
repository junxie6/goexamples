package main

import (
	"log"
	"net/http"
)

func srvHome(w http.ResponseWriter, r *http.Request) {
	log.Printf("URL %v", r.URL.Path)

	// The "/" pattern matches everything, so we need to check
	// that we're at the root here.
	// Reference:
	// https://golang.org/pkg/net/http/#example_ServeMux_Handle
	// https://golang.org/pkg/net/http/#ServeMux
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("home"))
}

func srvHello(w http.ResponseWriter, r *http.Request) {
	log.Printf("URL %v", r.URL.Path)

	w.Write([]byte("hello"))
}
func srvHelloWorld(w http.ResponseWriter, r *http.Request) {
	log.Printf("URL %v", r.URL.Path)

	w.Write([]byte("hello world"))
}

func main() {
	http.HandleFunc("/", srvHome)
	http.HandleFunc("/hello", srvHello)
	http.HandleFunc("/hello/world", srvHelloWorld)

	http.ListenAndServe(":8080", nil)
}
