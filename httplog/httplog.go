package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
)

/*
 * Log all http requestes and responses
 *
 * Reference:
 * https://play.golang.org/p/ND1HlS8dDn
 * https://groups.google.com/forum/#!topic/golang-nuts/s7Xk1q0LSU0
 * https://gist.github.com/JalfResi/87e8446b1bbf9506ce4143f89c7dfc9b
 */
func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)

		//handler.ServeHTTP(w, r)
		//return

		var dump []byte
		var err error

		// Dump request
		dump, err = httputil.DumpRequest(r, true)

		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}

		log.Printf("Dump Request:\n%s\n", dump)

		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, r)

		//log.Printf("Dump Response: %s", rec.Body)
		//return

		// Dump response
		dump, err = httputil.DumpResponse(rec.Result(), true)

		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}

		log.Printf("Dump Response:\n%s\n", dump)

		// we copy the captured response headers to our new response
		for k, v := range rec.Header() {
			w.Header()[k] = v
		}

		// grab the captured response code and response body
		w.WriteHeader(rec.Code)
		w.Write(rec.Body.Bytes())
	})
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	})

	http.ListenAndServe(":8080", Log(http.DefaultServeMux))
}
