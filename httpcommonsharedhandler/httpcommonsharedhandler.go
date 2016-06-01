package main

import (
	"github.com/NYTimes/gziphandler"
	"log"
	"net/http"
)

func serveTest1(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Test1"))
}

func serveTest2(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Test2"))
}

// Reference:
// http://stackoverflow.com/questions/26204485/gorilla-mux-custom-middleware
// http://laicos.com/writing-handsome-golang-middleware/
// https://www.nicolasmerouze.com/middlewares-golang-best-practices-examples/
// https://medium.com/@matryer/writing-middleware-in-golang-and-how-go-makes-it-so-much-fun-4375c1246e81#.4xnxwa6mx
func serveCommon(h http.Handler) http.Handler {
	return gziphandler.GzipHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("middleware: ", r.URL)

		h.ServeHTTP(w, r)
	}))
}

func main() {
	http.Handle("/test1", serveCommon(http.HandlerFunc(serveTest1)))
	http.Handle("/test2", serveCommon(http.HandlerFunc(serveTest2)))

	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
