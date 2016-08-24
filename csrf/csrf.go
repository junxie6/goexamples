package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	HTTP_PORT     = ":8080"
	CSRF_AUTH_KEY = "31-byte-long-auth-key----------"
	STATIC_DIR    = "/static/"
)

type RetJson struct {
	Status bool
}

type User struct {
	Uid   int
	Name  string
	Email string
}

func srvHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	w.Write([]byte(`<html>
	<head>
		<script src="https://code.jquery.com/jquery-2.2.4.min.js"></script>
		<script src="/static/debug.js"></script>
	</head>
	<body>
		X-CSRF-Token: <input type="text" id="CSRFToken" style="width: 800px;" />
		<br>Username: <input type="text" id="Username" value="jun" />
		<br>Password: <input type="text" id="Password" value="junpass" />
		<br>Response: <textarea id="JSONResponse" style="width: 300px;"></textarea>
		<br><button id="CSRFBtn">Get CSRF Token</button>
		<br><button id="LoginBtn">Login</button>
		<br><button id="SalOrderBtn">SalOrder</button>
	</body>
</html>
	`))
}

func srvNews(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func srvLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// We can trust that requests making it this far have satisfied
	// our CSRF protection requirements.

	o := struct {
		Status bool
	}{
		Status: true,
	}

	b, err := json.Marshal(o)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

func srvCSRFToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Get the token and pass it in the CSRF header. Our JSON-speaking client
	// or JavaScript framework can now read the header and return the token in
	// in its own "X-CSRF-Token" request header on the subsequent POST.
	w.Header().Set("X-CSRF-Token", csrf.Token(r))

	o := struct {
		Status bool
	}{
		Status: true,
	}

	b, err := json.Marshal(o)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

func srvSalOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	o := struct {
		Status bool
	}{
		Status: true,
	}

	b, err := json.Marshal(o)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

func main() {
	r := mux.NewRouter()

	r.PathPrefix(STATIC_DIR).Handler(http.StripPrefix(STATIC_DIR, http.FileServer(http.Dir("."+STATIC_DIR))))

	r.HandleFunc("/", srvHome)

	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/GetCSRFToken", srvCSRFToken).Methods("GET")
	api.HandleFunc("/News", srvNews).Methods("GET")
	api.HandleFunc("/SalOrder", srvSalOrder).Methods("POST")

	// All POST requests without a valid token will return HTTP 403 Forbidden.
	api.HandleFunc("/Login", srvLogin).Methods("POST")

	// Note: Don't forget to pass csrf.Secure(false) if you're developing locally over plain HTTP (just don't leave it on in production).
	err := http.ListenAndServe(HTTP_PORT, csrf.Protect([]byte(CSRF_AUTH_KEY), csrf.Secure(false), csrf.MaxAge(86400*1))(r))

	if err != nil {
		fmt.Printf("main(): %s\n", err)
	}
}
