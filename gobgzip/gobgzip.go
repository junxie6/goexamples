package main

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"io"
)

// Person ...
type Person struct {
	Name string
	Age  int
}

func main() {
	person := Person{
		Name: "AAA",
		Age:  123,
	}

	var encodedData []byte
	var gzippedData []byte
	var ungzippedData []byte
	var err error

	// encoding data to gob
	if encodedData, err = GobEncode(person); err != nil {
		fmt.Printf("err: %v\n", err)
	} else {
		fmt.Printf("encodedData: %v\n\n", encodedData)
	}

	// compressing data
	if gzippedData, err = GzipCompress(encodedData); err != nil {
		fmt.Printf("err: %v\n", err)
	} else {
		fmt.Printf("gzippedData: %v\n\n", gzippedData)
	}

	// uncompressing data
	if ungzippedData, err = GzipUncompress(bytes.NewReader(gzippedData)); err != nil {
		fmt.Printf("err: %v\n", err)
	} else {
		fmt.Printf("ungzippedData: %v\n\n", ungzippedData)
	}

	// decoding data from gob
	person2 := Person{}

	if err := GobDecode(bytes.NewReader(ungzippedData), &person2); err != nil {
		fmt.Printf("err: %v\n", err)
	} else {
		fmt.Printf("person2: %#v\n\n", person2)
	}

	fmt.Printf("Note: for some data, the compressed data is actually bigger than the original data.\n")
}

// GzipCompress ...
func GzipCompress(data []byte) ([]byte, error) {
	b := new(bytes.Buffer)
	gz := gzip.NewWriter(b)

	if _, err := gz.Write(data); err != nil {
		return nil, err
	}

	// Note: One more thing is that Flush() only writes the current data to the buffer.
	// It doesn't finish off the whole GZIP format.
	// So, in this case, it's pretty useless, since what's written on the last line is not a valid GZIP structure.
	// You need to call Close() before you do anything with the buffer.
	//
	// You can just call Close(). From the golang documentation: func (z *Writer) Close() error Close closes the Writer,
	// flushing any unwritten data to the underlying io.Writer, but does not close the underlying io.Writer.
	// http://golang.org/pkg/compress/gzip/#Writer.Close
	//
	// Just to add a bit, using defer to close the compressed Writer can lead to subtle bugs where
	// the buffer is being used before being closed.
	// This can result in unexpected EOF errors when reading the compressed data. Watch out!
	// Reference:
	// http://stackoverflow.com/questions/19197874/how-can-i-use-gzip-on-a-string-in-golang
	if err := gz.Flush(); err != nil {
		return nil, err
	}

	if err := gz.Close(); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// GzipUncompress ...
func GzipUncompress(data *bytes.Reader) ([]byte, error) {
	gz, _ := gzip.NewReader(data)
	gz.Close()

	// to standard output
	network := new(bytes.Buffer)
	_, _ = io.Copy(network, gz)
	return network.Bytes(), nil
}

// GobEncode ...
func GobEncode(obj interface{}) ([]byte, error) {
	// Initialize the encoder and decoder.  Normally enc and dec would be
	// bound to network connections and the encoder and decoder would
	// run in different processes.
	//var network bytes.Buffer        // Stand-in for a network connection
	network := new(bytes.Buffer)
	enc := gob.NewEncoder(network) // Will write to network.

	// Encode (send) the value.
	if err := enc.Encode(obj); err != nil {
		return nil, err
	}

	return network.Bytes(), nil
}

// GobDecode ...
func GobDecode(data *bytes.Reader, obj interface{}) error {
	dec := gob.NewDecoder(data) // Will read from network.

	// Decode (receive) the value.
	if err := dec.Decode(obj); err != nil {
		return err
	}
	return nil
}
