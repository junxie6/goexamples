package hello

import (
	"log"
	"net/http"
)

// SrvHello ...
func SrvHello(w http.ResponseWriter, r *http.Request) {
	log.Printf("act 1: %v\n", r.FormValue("act"))

	switch r.FormValue("act") {
	case "say":
		SayHello(r)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte(`{"Status":true,"ErrArr":null,"ErrCount":0,"ObjArr":null,"Data":{"name":"bot"}}`))
}

// SayHello ...
func SayHello(r *http.Request) {
	// DEBUG:
	log.Printf("act 2: %v\n", r.FormValue("act"))
	log.Printf("request body: %s\n", r.Body)
}
