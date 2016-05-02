package main

import (
	"fmt"
	"log"
	"net/http"
)

func serveExample1(w http.ResponseWriter, r *http.Request) {
	str := "Hello World"

	fmt.Fprintf(w, "URL: %s\n\nstr: %s\n\n", r.URL.Path[1:], str)

	w.Write([]byte(str))
}

func main() {

	http.HandleFunc("/ex1", serveExample1)

	err := http.ListenAndServe(":80", nil)

	if err != nil {
		log.Println(err)
		log.Printf("%v", err)
		log.Fatal(err)

		fmt.Printf("%v", err)
		fmt.Println(err)
		fmt.Errorf("%v", err)
	}
}
