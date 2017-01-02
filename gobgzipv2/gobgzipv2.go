package main

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"io"
	"log"
)

// Person ...
type Person struct {
	Name   string
	Age    int
	Money  float64
	Houses []House
}

// House ...
type House struct {
	StreetNum     int
	StreetName    string
	PropertyValue float64
}

func main() {
	person := GetPerson()
	person2 := Person{}

	if b, err := EncodeGobThenGzip(&person); err != nil {
		fmt.Printf("err: %v\n\n", err)
	} else if err := UngzipThenDecodeGob(bytes.NewReader(b), &person2); err != nil {
		fmt.Printf("err: %v\n\n", err)
	} else {
		fmt.Printf("gob encoded: %v\n\n", b)
		fmt.Printf("person2: %#v\n\n", person2)
	}
}

func GetPerson() Person {
	return Person{
		Name:  "AAA",
		Age:   123,
		Money: 123456.789,
		Houses: []House{
			House{
				StreetNum:     111,
				StreetName:    "Hello 111",
				PropertyValue: 123456.789,
			},
			House{
				StreetNum:     222,
				StreetName:    "Hello 222",
				PropertyValue: 987654.321,
			},
		},
	}
}

// EncodeGobThenGzip ...
func EncodeGobThenGzip(obj interface{}) ([]byte, error) {
	var gz *gzip.Writer
	var err error
	b := new(bytes.Buffer)

	if gz, err = gzip.NewWriterLevel(b, gzip.DefaultCompression); err != nil {
		return nil, err
	}

	if err := gob.NewEncoder(gz).Encode(obj); err != nil {
		return nil, err
	}

	// Note: gzip Writer must be closed before the bufer is being used.
	// using defer to close the compressed Writer can lead to subtle bugs where the buffer is being used before being closed.
	// This can result in unexpected EOF errors when reading the compressed data. Watch out!
	//
	// Reference:
	// http://stackoverflow.com/questions/19197874/how-can-i-use-gzip-on-a-string-in-golang
	// https://www.datadoghq.com/blog/crossing-streams-love-letter-gos-io-reader/
	if err := gz.Close(); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// UngzipThenDecodeGob ...
func UngzipThenDecodeGob(data io.Reader, obj interface{}) error {
	var gz *gzip.Reader
	var err error

	if gz, err = gzip.NewReader(data); err != nil {
		return err
	}

	defer func() {
		if err := gz.Close(); err != nil {
			log.Printf("err: %v\n", err)
		}
	}()

	if err := gob.NewDecoder(gz).Decode(obj); err != nil {
		return err
	}

	return nil
}
