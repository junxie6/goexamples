package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

const (
	HTTP_PORT     = ":80"
	CSRF_AUTH_KEY = "31-byte-long-auth-key"
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

func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	w.Write([]byte("<script src='/static/scripts/jquery-2.2.2.min.js'></script>" +
		"<script src='/static/scripts/test_form.js'></script>"))
}

func SubmitSignupForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// We can trust that requests making it this far have satisfied
	// our CSRF protection requirements.

	retjson := RetJson{}
	retjson.Status = true

	b, err := json.Marshal(retjson)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(b)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	uid := mux.Vars(r)["uid"] // variable name is case sensitive

	uid2, err := strconv.Atoi(uid)

	fmt.Printf("main(): %s\n", err)

	// Authenticate the request, get the id from the route params,
	// and fetch the user from the DB, etc.
	user := User{Uid: uid2, Name: "Test Name", Email: "test@abc.com"}

	// Get the token and pass it in the CSRF header. Our JSON-speaking client
	// or JavaScript framework can now read the header and return the token in
	// in its own "X-CSRF-Token" request header on the subsequent POST.
	w.Header().Set("X-CSRF-Token", csrf.Token(r))

	b, err := json.Marshal(user)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(b)
}

func main() {
	r := mux.NewRouter()

	r.PathPrefix(STATIC_DIR).Handler(http.StripPrefix(STATIC_DIR, http.FileServer(http.Dir("."+STATIC_DIR))))

	r.HandleFunc("/", RootHandler)

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/user/{uid}", GetUser).Methods("GET")

	// All POST requests without a valid token will return HTTP 403 Forbidden.
	api.HandleFunc("/signup/post", SubmitSignupForm).Methods("POST")

	// Note: Don't forget to pass csrf.Secure(false) if you're developing locally over plain HTTP (just don't leave it on in production).
	err := http.ListenAndServe(HTTP_PORT, csrf.Protect([]byte(CSRF_AUTH_KEY), csrf.Secure(false))(r))

	if err != nil {
		fmt.Printf("main(): %s\n", err)
	}
}
