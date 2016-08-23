package main

import (
	"log"
	"net/http"
	"time"
)

const (
	srvAddr           = "0.0.0.0:8080"
	srvWriteTimeout   = 15 * time.Second
	srvReadTimeout    = 15 * time.Second
	srvMaxHeaderBytes = 1 << 20
)

func srvHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

func main() {
	// Reference:
	// https://golang.org/pkg/net/http/#ServeMux
	// https://golang.org/pkg/net/http/#example_ServeMux_Handle
	// http://golang.org/pkg/net/http/#Server
	mux := http.NewServeMux()

	mux.HandleFunc("/", srvHome)

	//
	srv := &http.Server{
		Handler: mux,
		Addr:    srvAddr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout:   srvWriteTimeout,
		ReadTimeout:    srvReadTimeout,
		MaxHeaderBytes: srvMaxHeaderBytes,
	}

	log.Printf("Listening to TCP address: %v", srv.Addr)

	//
	if err := srv.ListenAndServe(); err != nil {
		log.Printf("srv.ListenAndServe: %v", err.Error())
	}
}
