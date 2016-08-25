package main

import (
	"log"
	"net/http"
)

type (
//MyHandler func(http.ResponseWriter, *http.Request, *ioxer.IOXer)
)

func srvHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Home\n"))
}

func something() {
}

func HelloHandler(h1 http.Handler) http.Handler {
	return func(h2 http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello Before\n"))
			h2.ServeHTTP(w, r)
			w.Write([]byte("Hello After\n"))
		})
	}(h1)
}

func WorldHandler(h1 http.Handler) http.Handler {
	return func(h2 http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("World Before\n"))
			h2.ServeHTTP(w, r)
			w.Write([]byte("World After\n"))
		})
	}(h1)
}

func ThereHandler(h1 http.Handler) http.Handler {
	return func(h2 http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("There Before\n"))
			h2.ServeHTTP(w, r)
			w.Write([]byte("There After\n"))
		})
	}(h1)
}

func s1(s string) string {
	log.Printf("s1")
	return "s1"
}
func s2(s string) string {
	log.Printf("s2")
	return "s2"
}
func s3(s string) string {
	log.Printf("s3")
	return "s3"
}

func main() {
	http.Handle("/", HelloHandler(WorldHandler(ThereHandler(http.HandlerFunc(srvHome)))))
	http.ListenAndServe(":8080", nil)

	//log.Printf("Here: %v", s1(s2(s3("asdf"))))
}
