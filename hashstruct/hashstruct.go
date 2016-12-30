package main

import (
	//"crypto/sha1"
	"crypto/sha256"
	"fmt"
)

type Foo struct {
	A int
	B string
}

// Reference:
// http://stackoverflow.com/questions/16330490/in-go-how-can-i-convert-a-struct-to-a-byte-array
func main() {
	h256 := sha256.New()
	foo := Foo{A: 1, B: "bar1"}

	s := fmt.Sprintf("%v", foo)

	sum256_1 := sha256.Sum256([]byte(s))
	sum256_2 := h256.Sum([]byte(s))

	fmt.Printf("%s hashes to %x\n", s, sum256_1)
	fmt.Printf("%s hashes to %x\n", s, sum256_2)
}
