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
// https://www.jokecamp.com/blog/examples-of-creating-base64-hashes-using-hmac-sha256-in-different-languages/
// https://stackoverflow.com/questions/18492576/share-cookie-between-subdomain-and-domain
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
