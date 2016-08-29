package main

import (
	"log"
	"net/http"
)

type (
//MyHandler func(http.ResponseWriter, *http.Request, *iojson.IOJSON)
)

func srvHome(w http.ResponseWriter, r *http.Request) {
	log.Printf("Home")
	w.Write([]byte("Home " + r.FormValue("act") + "\n"))
}

func something() {
}

//
func HelloHandler(h1 http.Handler) http.Handler {
	log.Printf("Inside Hello")
	return func(h2 http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Hello Before")
			w.Write([]byte("Hello Before " + r.FormValue("act") + "\n"))

			h2.ServeHTTP(w, r)

			log.Printf("Hello After")
			w.Write([]byte("Hello After " + r.FormValue("act") + "\n"))
		})
	}(h1)
}

func WorldHandler(h1 http.Handler) http.Handler {
	log.Printf("Inside World")
	return func(h2 http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("World Before")
			w.Write([]byte("World Before " + r.FormValue("act") + "\n"))

			h2.ServeHTTP(w, r)

			w.Write([]byte("World After " + r.FormValue("act") + "\n"))
			log.Printf("World After")
		})
	}(h1)
}

func ThereHandler(h1 http.Handler) http.Handler {
	log.Printf("Inside There")
	return func(h2 http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("There Before")
			w.Write([]byte("There Before " + r.FormValue("act") + "\n"))

			h2.ServeHTTP(w, r)

			log.Printf("There After")
			w.Write([]byte("There After " + r.FormValue("act") + "\n"))
		})
	}(h1)
}

//

//
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
	http.Handle("/", http.HandlerFunc(srvHome))
	http.Handle("/Hello", HelloHandler(WorldHandler(ThereHandler(http.HandlerFunc(srvHome)))))
	http.Handle("/How", http.HandlerFunc(srvHome))
	http.ListenAndServe(":8080", nil)

	//log.Printf("Here: %v", s1(s2(s3("asdf"))))
}
