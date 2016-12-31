package main

import (
	"bytes"
	"testing"
)

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
