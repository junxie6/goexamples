package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Host: %v", r.Host)
	// r.URL.Scheme will be empty if you're accessing the HTTP server not from an HTTP proxy,
	// a browser can issue a relative HTTP request instead of a absolute URL.
	// Additionally, you could check in the server/handler whether you get a
	// relative or absolute URL in the request by calling the IsAbs() method.
	// Reference: http://stackoverflow.com/questions/6899069/why-are-request-url-host-and-scheme-blank-in-the-development-server
	log.Printf("Scheme: %v", r.URL.Scheme)
	log.Printf("IsAbs: %v", r.URL.IsAbs())

	if _, err := w.Write([]byte("Gorilla!\n")); err != nil {
	}
}

func HmmHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Hmm!\n")); err != nil {
	}
}

func SalOrderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	IDSalOrder := vars["IDSalOrder"]

	if _, err := w.Write([]byte("SalOrder: " + IDSalOrder + "\n")); err != nil {
	}
}

func shared() {
}

func main() {
	r := mux.NewRouter()

	// Only matches if domain is "www.example.com".
	s := r.Host("erp.local").
		MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
			log.Printf("act: %v", r.FormValue("act"))

			if r.FormValue("act") == "test" {
				log.Printf("It is TRUE")
				return true
			} else {
				log.Printf("It is FALSE")
				return false
			}
		}).
		Subrouter()
	s.HandleFunc("/SalOrder/{IDSalOrder}", SalOrderHandler)
	s.HandleFunc("/SalOrder", SalOrderHandler)

	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/hmm", HmmHandler)

	//http.Handle("/", s)

	if err := http.ListenAndServe(":8080", r); err != nil {
	}
}
