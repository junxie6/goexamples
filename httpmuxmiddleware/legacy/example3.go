package main

import (
	"github.com/junxie6/alice"
	"github.com/junxie6/iojson"
	"log"
	"net/http"
	"time"
)

const (
	//sessionSecret     = "something-very-secret"
	//sessionName       = "MySession"
	srvAddr           = "0.0.0.0:8443"
	srvWriteTimeout   = 35 * time.Second
	srvReadTimeout    = 35 * time.Second
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

func LoggerHandler(h http.Handler) http.Handler {
	return NewLogger(h)
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//time.Sleep(5 * time.Second)

	//log.Printf("Before: %v\t%v", r.Host, r.URL.Path)

	l.handler.ServeHTTP(w, r)

	//log.Printf("After: %v\t%v", r.Host, r.URL.Path)
}

//
func mySingleHost(h http.Handler) http.Handler {
	return SingleHost(h, "erp.local:8443")
}

func SingleHost(handler http.Handler, allowedHost string) http.Handler {
	ourFunc := func(w http.ResponseWriter, r *http.Request) {
		o := r.Context().Value("iojson").(*iojson.IOJSON)

		host := r.Host

		if host == allowedHost {
			o.AddData("SingleHost", "Inside SingleHost")

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

func srvHome2(w http.ResponseWriter, r *http.Request) {
	//log.Printf("Inside home2")

	o := r.Context().Value("iojson").(*iojson.IOJSON)
	o.AddData("act", "home 2; "+r.FormValue("act"))
}

func srvHome3(w http.ResponseWriter, r *http.Request) {
	//log.Printf("Inside home3")

	o := r.Context().Value("iojson").(*iojson.IOJSON)
	o.AddData("act", "home 3; "+r.FormValue("act"))
}

func srvNews2(w http.ResponseWriter, r *http.Request) {
	o := iojson.NewIOJSON()
	o.AddData("SingleHost", "Inside SingleHost")
	o.AddData("act", r.FormValue("act"))
	o.Echo(w)
}

func main() {
	chain := alice.New(
		iojson.EchoHandler,
		LoggerHandler,
		mySingleHost,
	)

	mux := http.NewServeMux()

	mux.Handle("/Home2", chain.Then(http.HandlerFunc(srvHome2)))
	mux.Handle("/Home3", chain.Then(http.HandlerFunc(srvHome3)))

	mux.Handle("/News2", http.HandlerFunc(srvNews2))

	srv := &http.Server{
		Handler: mux,
		Addr:    srvAddr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout:   srvWriteTimeout,
		ReadTimeout:    srvReadTimeout,
		MaxHeaderBytes: srvMaxHeaderBytes,
	}

	//if err := srv.ListenAndServe(); err != nil {
	if err := srv.ListenAndServeTLS("mydomain.com.crt", "mydomain.com.key"); err != nil {
		log.Printf("srv.ListenAndServe: %v", err.Error())
	}
}
