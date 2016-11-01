package main

import (
	"log"
	"net/http"
	"strings"
	"time"
)

import (
	"github.com/julienschmidt/httprouter"
)

const (
	srvMaxHeaderBytes = 1 << 20
	gStaticDir        = "static"
)

var (
	gMux *httprouter.Router
)

func main() {
	gMux = httprouter.New()

	// Serving the static files - basic
	//gMux.ServeFiles("/"+gStaticDir+"/*filepath", http.Dir(gStaticDir))

	// Serving the static files - use this if you would like to work with some middlewares.
	// NOTE: FileServer calls path.Clean() to clean up path.
	gMux.Handler("GET", "/"+gStaticDir+"/*filepath", NoDirListingHandler(http.StripPrefix("/"+gStaticDir+"/", http.FileServer(http.Dir(gStaticDir)))))

	srv := &http.Server{
		Handler:        gMux,
		Addr:           ":8080",
		WriteTimeout:   time.Duration(20) * time.Second, // Good practice: enforce timeouts for servers you create.
		ReadTimeout:    time.Duration(20) * time.Second,
		MaxHeaderBytes: srvMaxHeaderBytes,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("http.ListenAndServe: ", err)
	}
}

func NoDirListingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		h.ServeHTTP(w, r)
	})
}
