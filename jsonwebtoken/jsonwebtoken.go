package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"text/template"
	"time"
)

const (
	Bearer = "Bearer"

	homeHTML = `<!DOCTYPE html>
<html lang="en">
    <head>
        <title>Template Example</title>
    </head>
    <body>
        {{.Data}}
    </body>
</html>
`
)

var (
	homeTempl = template.Must(template.New("").Parse(homeHTML))
)

func mykeys(key string) (string, error) {
	strMap := map[string]string{
		"username1": "mysecret1",
		"username2": "mysecret2",
	}

	if val, ok := strMap[key]; ok {
		return val, nil
	}

	return "", fmt.Errorf("MY: key not found")
}

func genJsonWebToken(username string, signingKey string) string {
	// New web token.
	token := jwt.New(jwt.SigningMethodHS256)

	// Set a header and a claim
	token.Header["typ"] = "JWT"

	token.Claims["username"] = username
	token.Claims["exp"] = time.Now().Add(time.Hour * 96).Unix()
	token.Claims["foo"] = "bar"

	// Generate encoded token
	t, _ := token.SignedString([]byte(signingKey))

	return t
}

func serveValidateJsonWebToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	auth := r.Header.Get("Authorization")
	l := len(Bearer)

	//str := r.RemoteAddr

	fmt.Printf("Here 1: %v\n", auth)

	if len(auth) > l+1 && auth[:l] == Bearer {
		fmt.Printf("Here 2: %v\n", auth[l+1:])

		t, err := jwt.Parse(auth[l+1:], func(token *jwt.Token) (interface{}, error) {

			// Always check the signing method
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				//fmt.Printf("Here 3: %v\n", "not good")
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// Return the key for validation
			fmt.Printf("Here username: %v\n", token.Header["username"])
			fmt.Printf("Here token: %v\n", token)

			username := token.Header["username"].(string)

			secretkey, err := mykeys(username)

			if err == nil {
				return []byte(secretkey), nil
			} else {
				return []byte(secretkey), nil
			}
		})

		fmt.Printf("Here 5: %v\n", t)

		if err == nil && t.Valid {
			w.Write([]byte("Good " + t.Header["username"].(string)))
		} else {
			w.Write([]byte("Bad"))
		}
	} else {
		w.Write([]byte("<script src='/static/scripts/jquery-2.2.2.min.js'></script>"))
	}
}

func servejwt(w http.ResponseWriter, r *http.Request) {
	testMap := map[string]interface{}{
		"status": false,
	}

	token := genJsonWebToken("test", "test")
	testMap["token"] = token

	testJSON, err := json.Marshal(testMap)

	if err != nil {
		w.Write([]byte("{\"status\": false}"))
	}

	w.Write([]byte(testJSON))
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var v = struct {
		Host string
		Data string
	}{
		r.Host,
		"Test",
	}

	homeTempl.Execute(w, &v)
}

func main() {

	http.HandleFunc("/validatejwt", serveValidateJsonWebToken)
	http.HandleFunc("/jwt", servejwt)

	http.HandleFunc("/", serveHome)

	err := http.ListenAndServe(":80", nil)

	if err != nil {
		fmt.Printf("main(): %s\n", err)
		log.Fatal("ListenAndServe: ", err)
	}
}
