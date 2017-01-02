package main

import (
	"bytes"
	"fmt"
	"testing"
)

func BenchmarkGobGzip(b *testing.B) {
	for i := 0; i < b.N; i++ {
		person := GetPerson()
		person2 := Person{}

		if b, err := GobEncode(&person); err != nil {
			fmt.Printf("err: %v\n", err)
		} else if b2, err := GzipCompress(b); err != nil {
			fmt.Printf("err: %v\n", err)
		} else if b3, err := GzipUncompress(bytes.NewReader(b2)); err != nil {
			fmt.Printf("err: %v\n", err)
		} else if err := GobDecode(bytes.NewReader(b3), &person2); err != nil {
			fmt.Printf("err: %v\n", err)
		}
	}
}

func BenchmarkByteNewBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bytes.NewBuffer([]byte("R29waGVycyBydWxlIQ=="))
	}
}

func BenchmarkByteNewReader(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bytes.NewReader([]byte("R29waGVycyBydWxlIQ=="))
	}
}
