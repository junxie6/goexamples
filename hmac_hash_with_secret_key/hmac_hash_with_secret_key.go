package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

// Reference:
// https://www.calhoun.io/securing-cookies-in-go/
// https://github.com/danharper/hmac-examples
func main() {
	secret := []byte("the shared secret key here")
	message := []byte("the message to hash here")

	hash := hmac.New(sha256.New, secret)
	hash.Write(message)

	// to lowercase hexits
	s1 := hex.EncodeToString(hash.Sum(nil))

	// to base64
	s2 := base64.StdEncoding.EncodeToString(hash.Sum(nil))

	fmt.Printf("HERE: %v\n", s1)
	fmt.Printf("HERE: %v\n", s2)
}
