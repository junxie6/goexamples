package main

import (
	"github.com/NYTimes/gziphandler"
	"log"
	"net/http"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

func main() {
	http.Handle("/", gziphandler.GzipHandler(http.HandlerFunc(serveHome)))

	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
