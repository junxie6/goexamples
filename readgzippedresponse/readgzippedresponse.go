package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Reference:
// https://www.socketloop.com/tutorials/golang-read-gzipped-http-response
func main() {
	client := new(http.Client)

	request, err := http.NewRequest("Get", "http://blog.ijun.org/", nil)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	request.Header.Add("Accept-Encoding", "gzip")

	response, err := client.Do(request)

	if response != nil {
		defer response.Body.Close()
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Check that the server actual sent compressed data
	var reader io.ReadCloser

	fmt.Printf("Content-Encoding: %v\n", response.Header.Get("Content-Encoding"))

	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer reader.Close()
	default:
		reader = response.Body
	}

	// to standard output
	_, err = io.Copy(os.Stdout, reader)

	// see https://www.socketloop.com/tutorials/golang-saving-and-reading-file-with-gob
	// on how to save to file

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
