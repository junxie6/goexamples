package main

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
)

// Person ...
type Person struct {
	Name   string
	Age    int
	Money  float64
	Houses []House
}

type House struct {
	StreetNum     int
	StreetName    string
	PropertyValue float64
}

func main() {
	person := Person{
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

	b, _ := EncodeGobGzip(&person)

	person2 := Person{}

	UngzipDecodeGob(bytes.NewReader(b), &person2)

	fmt.Printf("gob encoded: %v\n\n", b)
	fmt.Printf("person2: %#v\n\n", person2)
}

func EncodeGobGzip(obj interface{}) ([]byte, error) {
	b := new(bytes.Buffer)

	enc := gob.NewEncoder(b)

	if err := enc.Encode(obj); err != nil {
		return nil, err
	}

	//
	var gz *gzip.Writer
	var err error
	b2 := new(bytes.Buffer)

	if gz, err = gzip.NewWriterLevel(b2, gzip.DefaultCompression); err != nil {
		return nil, err
	}

	if _, err := gz.Write(b.Bytes()); err != nil {
		return nil, err
	}

	if err := gz.Close(); err != nil {
		return nil, err
	}

	return b2.Bytes(), nil
}

func UngzipDecodeGob(data io.Reader, obj interface{}) error {
	//b := new(bytes.Buffer)
	var gz *gzip.Reader
	var err error

	if gz, err = gzip.NewReader(data); err != nil {
		return err
	}

	if err := gz.Close(); err != nil {
		return err
	}

	dec := gob.NewDecoder(gz) // Will read from network.

	// Decode (receive) the value.
	if err := dec.Decode(obj); err != nil {
		return err
	}

	return nil
}

func LoadGzippedJSON(r io.Reader, v interface{}) error {
	raw, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	return json.NewDecoder(raw).Decode(&v)
}
