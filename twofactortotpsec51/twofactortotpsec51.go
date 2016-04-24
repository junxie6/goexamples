package main

import (
	"crypto"
	"fmt"
	"github.com/sec51/twofactor"
	"gopkg.in/redis.v3"
	"log"
	"net/http"
	"strconv"
)

var redisClient = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home: %s!", r.URL.Path[1:])
}

func writeImage(w http.ResponseWriter, r *http.Request) {
	otp, err := twofactor.NewTOTP("info@sec51.com", "Sec51", crypto.SHA1, 6)

	if err != nil {
		log.Println("call NewTOTP error.")
	}

	qrBytes, err := otp.QR()

	if err != nil {
		log.Println("call QR error.")
	}

	if b, err := otp.ToBytes(); err != nil {
		log.Println("call ToBytes error.")
	} else {
		log.Printf("Before: %v\n", b) // show the byte array before storing in Redis

		if err := redisClient.Set("info@sec51.com", b, 0).Err(); err != nil {
			log.Printf("H: %v\n\n", err)
		}
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(qrBytes)))

	if _, err := w.Write(qrBytes); err != nil {
		log.Println("unable to write image.")
	}
}

func verifyIt(w http.ResponseWriter, r *http.Request) {
	val, err := redisClient.Get("info@sec51.com").Bytes()

	log.Printf("Afterr: %v\n", val) // show the byte array after getting it from Redis

	otp, err := twofactor.TOTPFromBytes(val, "Sec51")

	if err != nil {
		log.Println("call TOTPFromBytes error.")
	}

	if err := otp.Validate(r.URL.Path[3:]); err != nil {
		fmt.Fprintf(w, "Wrong code: %s. Err: %v\n", r.URL.Path[3:], err)
	} else {
		fmt.Fprintf(w, "Good Job: %s.\n", r.URL.Path[3:])
	}
}

func main() {

	http.HandleFunc("/v/", verifyIt)
	http.HandleFunc("/qr/", writeImage)
	http.HandleFunc("/", handler)

	http.ListenAndServe(":80", nil)
}
