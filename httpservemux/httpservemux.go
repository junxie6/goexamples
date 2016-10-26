package main

import (
	"log"
	"net/http"
	"time"
)

const (
	srvMaxHeaderBytes = 1 << 20
)

var gMux *http.ServeMux

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World! " + r.URL.Path))
}

func main() {
	gMux = http.NewServeMux()

	gMux.Handle("/hello1", http.HandlerFunc(hello))
	gMux.HandleFunc("/hello2", hello)

	srv := &http.Server{
		Handler:        gMux,
		Addr:           ":8080",
		WriteTimeout:   time.Duration(15) * time.Second, // Good practice: enforce timeouts for servers you create.
		ReadTimeout:    time.Duration(15) * time.Second,
		MaxHeaderBytes: srvMaxHeaderBytes,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("http.ListenAndServe: ", err)
	}
}
