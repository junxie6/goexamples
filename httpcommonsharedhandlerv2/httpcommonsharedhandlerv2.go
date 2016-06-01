package main

import (
	"github.com/NYTimes/gziphandler"
	"log"
	"net/http"
)

// Reference:
// http://stackoverflow.com/questions/26204485/gorilla-mux-custom-middleware
// http://laicos.com/writing-handsome-golang-middleware/
// https://www.nicolasmerouze.com/middlewares-golang-best-practices-examples/
// https://medium.com/@matryer/writing-middleware-in-golang-and-how-go-makes-it-so-much-fun-4375c1246e81#.4xnxwa6mx

type (
	MyHandler func(w http.ResponseWriter, r *http.Request) error
)

func serveCommon(myhandlers ...MyHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, myhandler := range myhandlers {
			if err := myhandler(w, r); err != nil {
				w.Write([]byte(err.Error()))
				return
			}
		}
	})
}

func serveCheckAllowHost(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("Test1"))
	return nil
}

func serveTest1(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("Test1"))
	return nil
}

func main() {
	http.Handle("/test1", serveCommon(serveTest1))

	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
