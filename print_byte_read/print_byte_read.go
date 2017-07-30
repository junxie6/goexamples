// Display each byte read while reading the data
package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	//"os"
	//"strings"

	"github.com/davecgh/go-spew/spew"
)

// UseStdLibExample uses Go's stdlib io.TeeReader
// https://golang.org/pkg/io/#TeeReader
func UseStdLibExample() {
	var src io.Reader // Source file/url/etc
	var debugBuf bytes.Buffer

	// Create some random input data.
	src = bytes.NewBufferString("Some random input data")

	//
	r := io.TeeReader(src, &debugBuf)

	dst, err := ioutil.ReadAll(r)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n\n", dst)

	spew.Dump(debugBuf.Bytes())
}

func main() {

	// Use Go's stdlib io.TeeReader
	UseStdLibExample()
}
