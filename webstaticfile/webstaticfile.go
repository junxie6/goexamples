package main

import (
	"log"
	"net/http"
)

func main() {
	http.Handle("/static/", http.FileServer(http.Dir(".")))

	// For the serving directory name is different than the URL.
	// Note: Without http.StripPrefix function, if user accesses http://example.com/static/home.html, the server would look for it in ./public/static/home.html
	//http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./public"))))

	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
