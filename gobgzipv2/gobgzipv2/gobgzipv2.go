package gobgzipv2

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"io"
	"log"
)

// EncodeGobThenGzip ...
func EncodeGobThenGzip(obj interface{}) (io.Reader, error) {
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
	// https://medium.com/@matryer/golang-advent-calendar-day-seventeen-io-reader-in-depth-6f744bb4320b#.erh651rjt
	if err := gz.Close(); err != nil {
		return nil, err
	}

	return b, nil
}

// UngzipThenDecodeGob ...
func UngzipThenDecodeGob(r io.Reader, obj interface{}) error {
	var gz *gzip.Reader
	var err error

	if gz, err = gzip.NewReader(r); err != nil {
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

// EncodeGob ...
func EncodeGob(obj interface{}) (io.Reader, error) {
	b := new(bytes.Buffer)

	if err := gob.NewEncoder(b).Encode(obj); err != nil {
		return nil, err
	}

	return b, nil
}

// DecodeGob ...
func DecodeGob(r io.Reader, obj interface{}) error {
	if err := gob.NewDecoder(r).Decode(obj); err != nil {
		return err
	}

	return nil
}
