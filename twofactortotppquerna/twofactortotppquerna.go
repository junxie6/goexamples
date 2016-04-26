package main

import (
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"gopkg.in/redis.v3"

	"bufio"
	"bytes"
	"fmt"
	"image/png"
	//"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

var redisClient = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

func display(key *otp.Key, data []byte) {
	fmt.Printf("Issuer:       %s\n", key.Issuer())
	fmt.Printf("Account Name: %s\n", key.AccountName())
	fmt.Printf("Secret:       %s\n", key.Secret())
	fmt.Println("Writing PNG to qr-code.png....")
	//ioutil.WriteFile("qr-code.png", data, 0644)
	fmt.Println("")
	fmt.Println("Please add your TOTP to your OTP Application now!")
	fmt.Println("")
}

func promptForPasscode() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Passcode: ")
	text, _ := reader.ReadString('\n')
	return text
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home: %s!", r.URL.Path[1:])
}

func writeImage(w http.ResponseWriter, r *http.Request) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Example.com",
		AccountName: "test@example.com",
	})

	if err != nil {
		panic(err)
	}

	// Convert TOTP key into a PNG
	var buf bytes.Buffer

	img, err := key.Image(200, 200)

	if err != nil {
		panic(err)
	}

	png.Encode(&buf, img)

	// display the QR code to the user.
	qrBytes := buf.Bytes()
	display(key, qrBytes)

	if err := redisClient.Set("test@example.com", key.Secret(), 0).Err(); err != nil {
		log.Printf("H: %v\n\n", err)
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(qrBytes)))

	if _, err := w.Write(qrBytes); err != nil {
		log.Println("unable to write image.")
	}
}

func verifyIt(w http.ResponseWriter, r *http.Request) {
	// Now Validate that the user's successfully added the passcode.
	fmt.Println("Validating TOTP...")
	//passcode := promptForPasscode()

	val, err := redisClient.Get("test@example.com").Result()

	if err != nil {
		log.Println("call Result error.")
	}

	valid := totp.Validate(r.URL.Path[3:], val)

	if valid {
		fmt.Fprintf(w, "Good Job: %s.\n", r.URL.Path[3:])
	} else {
		fmt.Fprintf(w, "Wrong code: %s.\n", r.URL.Path[3:])
	}
}

func main() {

	http.HandleFunc("/v/", verifyIt)
	http.HandleFunc("/qr/", writeImage)
	http.HandleFunc("/", handler)

	http.ListenAndServe(":80", nil)
}
