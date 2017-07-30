// Display each byte read while reading the data
// Reference:
// https://stackoverflow.com/questions/22421375/how-to-print-the-bytes-while-the-file-is-being-downloaded-golang
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
	src = io.TeeReader(src, &debugBuf)

	ReadBytesAndPrint(src, &debugBuf)
}

func ReadBytesAndPrint(rd io.Reader, debugBuf *bytes.Buffer) {
	b, err := ioutil.ReadAll(rd)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n\n", b)
	spew.Dump((*debugBuf).Bytes())
}

// PrintRead wraps an existing io.Reader.
//
// It simply forwards the Read() call, while displaying
// the results from individual calls to it.
type PrintRead struct {
	reader   io.Reader
	debugBuf bytes.Buffer
	total    int64 // Total # of bytes transferred
}

// Read 'overrides' the underlying io.Reader's Read method.
// This is the one that will be called by io.Copy(). We simply
// use it to keep track of byte counts and then forward the call.
func (pr *PrintRead) Read(p []byte) (int, error) {
	n, err := pr.reader.Read(p)

	if n > 0 {
		pr.total += int64(n)

		if n, err := pr.debugBuf.Write(p[:n]); err != nil {
			return n, err
		}
	}

	return n, err
}

func UseCustomStruct() {
	var src io.Reader // Source file/url/etc

	// Create some random input data.
	src = bytes.NewBufferString("Some random input data")

	pr := &PrintRead{
		reader: src,
	}

	ReadBytesAndPrint(pr, &pr.debugBuf)

	fmt.Printf("Total bytes: %d\n", (*pr).total)
}

func main() {
	// Use Go's stdlib io.TeeReader
	UseStdLibExample()

	// Use a custom struct
	UseCustomStruct()
}
