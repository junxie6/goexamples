// Reference:
// https://medium.com/@xoen/golang-read-from-an-io-readwriter-without-loosing-its-content-2c6911805361
// https://github.com/gin-gonic/gin/issues/1295
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Person struct {
	Name string
}

func srvHome(w http.ResponseWriter, r *http.Request) {
	// Read the content
	//var bodyBytes []byte

	bodyBytes := make([]byte, 1024)
	numOfBytes := 0
	var err error

	if r.Body != nil {
		// http.Response.Body is of type io.ReadCloser, which can only be read once.
		// When you read from io.ReadCloser it drains it.
		// Once you read from it, the content is gone. You can’t read from it a second time.

		// Or use ioutil.ReadAll(r.Body)
		// var bodyBytes []byte
		//if bodyBytes, err = ioutil.ReadAll(r.Body); err != nil {
		if numOfBytes, err = r.Body.Read(bodyBytes); err != nil {
			if err != io.EOF {
				fmt.Fprintf(w, "Error: %s!", err.Error())
				return
			}
		}
	}

	//numOfBytes = len(bodyBytes)

	// Restore the io.ReadCloser to its original state
	//r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes[0:numOfBytes]))

	// Read the r.Body one more time
	p1 := Person{}

	if err = json.NewDecoder(r.Body).Decode(&p1); err != nil {
		fmt.Fprintf(w, "Error: %s!", err.Error())
		return
	}

	fmt.Fprintf(w, "Hello: %#v!", p1)
	//w.Write([]byte("Hello"))
}

func main() {

	http.HandleFunc("/", srvHome)

	http.ListenAndServe(":8080", nil)
}
