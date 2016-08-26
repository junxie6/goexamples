package main

import (
	"github.com/junhsieh/alice"
	"github.com/junhsieh/iojson"
	"log"
	"net/http"
	"time"
)

const (
	//sessionSecret     = "something-very-secret"
	//sessionName       = "MySession"
	srvAddr           = "0.0.0.0:8443"
	srvWriteTimeout   = 15 * time.Second
	srvReadTimeout    = 15 * time.Second
	srvMaxHeaderBytes = 1 << 20
	//CSRF_AUTH_KEY     = "31-byte-long-auth-key----------"
)

// Any struct with the method ServeHTTP(http.ResponseWriter, *http.Request) will be
// implementing http.Handler and will be usable with the Go muxer (http.Handle(pattern, handler) function).
type Logger struct {
	handler http.Handler
}

func NewLogger(handler http.Handler) *Logger {
	return &Logger{handler: handler}
}

func LoggerHandler(h http.Handler, o *iojson.IOJSON) http.Handler {
	return NewLogger(h)
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//time.Sleep(5 * time.Second)

	log.Printf("Before: %v\t%v", r.Host, r.URL.Path)

	l.handler.ServeHTTP(w, r)

	log.Printf("After: %v\t%v", r.Host, r.URL.Path)
}

func mySingleHost(h http.Handler, o *iojson.IOJSON) http.Handler {
	return SingleHost(h, o, "erp.local:8443")
}

func SingleHost(handler http.Handler, o *iojson.IOJSON, allowedHost string) http.Handler {
	ourFunc := func(w http.ResponseWriter, r *http.Request) {
		host := r.Host
		if host == allowedHost {
			handler.ServeHTTP(w, r)
		} else {
			//w.WriteHeader(403)
			//w.Write([]byte("Disallow host " + host + " " + r.FormValue("act") + "\n"))
			log.Printf("Disallow host " + host)
			o.AddError("Disallow host " + host + " " + r.FormValue("act") + "\n")
		}
	}
	return http.HandlerFunc(ourFunc)
}

func srvHome(h http.Handler, o *iojson.IOJSON) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Inside Home")
		o.AddData("ok", "okay la")
		o.AddData("ok", "okay la"+r.FormValue("act"))
	})
}

func main() {
	chain := alice.New(iojson.EchoHandler, LoggerHandler, mySingleHost, srvHome).Run()

	srv := &http.Server{
		Handler: chain,
		Addr:    srvAddr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout:   srvWriteTimeout,
		ReadTimeout:    srvReadTimeout,
		MaxHeaderBytes: srvMaxHeaderBytes,
	}

	//http.Handle("/", NewLogger(http.HandlerFunc(srvHome)))
	//http.HandleFunc("/", srvHome)
	//http.Handle("/mrtg/", http.StripPrefix("/mrtg/", http.FileServer(http.Dir("/var/www/html/mymrtg"))))

	//if err := srv.ListenAndServe(); err != nil {
	if err := srv.ListenAndServeTLS("mydomain.com.crt", "mydomain.com.key"); err != nil {
		log.Printf("srv.ListenAndServe: %v", err.Error())
	}
}
