package main

import (
	"fmt"
	"log"
	"net/http"
)

// Reference:
// https://play.golang.org/p/ND1HlS8dDn
// https://groups.google.com/forum/#!topic/golang-nuts/s7Xk1q0LSU0
func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	})
	http.ListenAndServe(":8080", Log(http.DefaultServeMux))
}
