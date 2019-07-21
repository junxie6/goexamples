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

func srvExample1(w http.ResponseWriter, r *http.Request) {
	// Read the content
	bodyBytes := make([]byte, 1024)
	numOfBytes := 0
	var err error

	if r.Body != nil {
		// http.Request.Body is of type io.ReadCloser, which can only be read once.
		// When you read from io.ReadCloser it drains it.
		// Once you read from it, the content is gone. You can’t read from it a second time.

		if numOfBytes, err = r.Body.Read(bodyBytes); err != nil {
			if err != io.EOF {
				fmt.Fprintf(w, "Error: %s!", err.Error())
				return
			}
		}
	}

	// One option is to read the body content from bodyBytes[0:numOfBytes]
	// or use ioutil.NopCloser to restore r.Body to its original state.
	p001 := Person{}

	if err = json.NewDecoder(bytes.NewReader(bodyBytes[0:numOfBytes])).Decode(&p001); err != nil {
		fmt.Fprintf(w, "Error: %s!", err.Error())
		return
	}

	p002 := Person{}

	if err = json.NewDecoder(bytes.NewReader(bodyBytes[0:numOfBytes])).Decode(&p002); err != nil {
		fmt.Fprintf(w, "Error: %s!", err.Error())
		return
	}

	// Restore the io.ReadCloser to its original state.
	// Use bytes.NewBuffer() to obtain an io.Reader from a byte slice.
	// Because bytes.Buffer does not have a Close() method.
	// For this you may use ioutil.NopCloser() which wraps an io.Reader,
	// and returns an io.ReadCloser, whose added Close() method will be a no-op (does nothing)
	// Reference: https://stackoverflow.com/questions/43021058/golang-read-request-body
	//
	// NOTE: or use bytes.NewReader() instead of bytes.NewBuffer()
	//r.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes[0:numOfBytes]))
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes[0:numOfBytes]))

	// Read the r.Body one more time
	p1 := Person{}

	if err = json.NewDecoder(r.Body).Decode(&p1); err != nil {
		fmt.Fprintf(w, "Error: %s!", err.Error())
		return
	}

	// Restore the io.ReadCloser to its original state
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes[0:numOfBytes]))

	// Read the r.Body one more time
	p2 := Person{}

	if err = json.NewDecoder(r.Body).Decode(&p2); err != nil {
		fmt.Fprintf(w, "Error: %s!", err.Error())
		return
	}

	fmt.Fprintf(w, "%#v\n%#v\n%#v\n%#v\n", p001, p002, p1, p2)
	//w.Write([]byte("Hello"))
}

func srvExample2(w http.ResponseWriter, r *http.Request) {
	var err error
	var bodyBytes []byte

	// NOTE: It's not a good idea to use ioutil.ReadAll when r.Body is very large, such as 1G Bytes. Use io.LimitedReader to prevent it.
	if bodyBytes, err = ioutil.ReadAll(r.Body); err != nil {
		fmt.Fprintf(w, "Error: %s!", err.Error())
		return
	}

	// Restore the io.ReadCloser to its original state
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	// Read the r.Body one more time
	p1 := Person{}

	if err = json.NewDecoder(r.Body).Decode(&p1); err != nil {
		fmt.Fprintf(w, "Error: %s!", err.Error())
		return
	}

	// Restore the io.ReadCloser to its original state
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	p2 := Person{}

	if err = json.NewDecoder(r.Body).Decode(&p2); err != nil {
		fmt.Fprintf(w, "Error: %s!", err.Error())
		return
	}

	fmt.Fprintf(w, "%#v\n%#v\n", p1, p2)
}
func srvExample3(w http.ResponseWriter, r *http.Request) {
	var err error

	// bytes.Buffer is an io.Reader and an io.Writer
	buf := bytes.NewBuffer(make([]byte, 0))

	// TeeReader returns a Reader that writes to w what it reads from r.
	reader := io.TeeReader(r.Body, buf)

	//
	p1 := Person{}

	// NOTE: We are using reader instead of r.Body
	if err := json.NewDecoder(reader).Decode(&p1); err != nil {
		fmt.Fprintf(w, "Error: %s!", err.Error())
		return
	}

	// Restore the io.ReadCloser to its original state
	r.Body = ioutil.NopCloser(bytes.NewBuffer(buf.Bytes()))

	// Read the r.Body one more time
	p2 := Person{}

	if err = json.NewDecoder(r.Body).Decode(&p2); err != nil {
		fmt.Fprintf(w, "Error: %s!", err.Error())
		return
	}

	fmt.Fprintf(w, "%#v\n%#v\n", p1, p2)
}

func srvExample4(w http.ResponseWriter, r *http.Request) {
	// bytes.Buffer is an io.Reader and an io.Writer
	buf := bytes.NewBuffer(make([]byte, 0))
	var numOfBytes int64
	var err error

	if numOfBytes, err = io.Copy(buf, r.Body); err != nil {
		fmt.Fprintf(w, "Error: %s!", err.Error())
		return
	}

	// Restore the io.ReadCloser to its original state
	r.Body = ioutil.NopCloser(bytes.NewBuffer(buf.Bytes()))

	p1 := Person{}

	if err := json.NewDecoder(r.Body).Decode(&p1); err != nil {
		fmt.Fprintf(w, "Error: %s!", err.Error())
		return
	}

	// Restore the io.ReadCloser to its original state
	r.Body = ioutil.NopCloser(bytes.NewBuffer(buf.Bytes()))

	// Read the r.Body one more time
	p2 := Person{}

	if err = json.NewDecoder(r.Body).Decode(&p2); err != nil {
		fmt.Fprintf(w, "Error: %s!", err.Error())
		return
	}

	fmt.Fprintf(w, "numOfBytes: %d\n%#v\n%#v\n", numOfBytes, p1, p2)
}

func main() {
	// Example1 uses r.Body.Read()
	// curl http://localhost:8080/Example1 --data '{"Name":"asdf2"}'
	http.HandleFunc("/Example1", srvExample1)

	// Example2 uses ioutil.ReadAll()
	// curl http://localhost:8080/Example2 --data '{"Name":"asdf2"}'
	http.HandleFunc("/Example2", srvExample2)

	// Example3 uses io.TeeReader()
	http.HandleFunc("/Example3", srvExample3)

	// Example4 uses io.Copy()
	http.HandleFunc("/Example4", srvExample4)

	http.ListenAndServe(":8080", nil)
}
