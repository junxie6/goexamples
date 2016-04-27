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
	bearer = "Bearer"

	homeHTML = `<!DOCTYPE html>
<html lang="en">
	<head>
		<title>Template Example</title>
		<script src="https://code.jquery.com/jquery-2.2.3.min.js"></script>
	  <script>
			$(document).ready(function() {
				var tokenJWT = "";

				$('#getJWTBtn').on('click', function(event){
					event.preventDefault();

					$.ajax({
						url: '/jwt',
						dataType: 'json',
					}).done(function(data) {
						$('#msg').html(JSON.stringify(data));

						if (data.status == true) {
							tokenJWT = data.token;
						}
					}).fail(function() {
					}).always(function() {
					});
					});

				$('#validateJWTBtn').on('click', function(event){
					event.preventDefault();

					$.ajax({
						url: '/validatejwt',
						dataType: 'json',
						headers: {
							"Authorization": "Bearer " + tokenJWT,
						},
					}).done(function(data) {
						$('#msg').html(JSON.stringify(data));
					}).fail(function() {
					}).always(function() {
					});
					});
			});
		</script>
	</head>
	<body>
		<p>{{.Data}}</p>

		<form>
			<button id="getJWTBtn">Get JWT</button>
			<button id="validateJWTBtn">Validate JWT</button>
			<p id="msg"></p>
		</form>
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

func serveValidateJWT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	outMap := map[string]interface{}{
		"status": false,
	}

	auth := r.Header.Get("Authorization")
	bearerLen := len(bearer)

	//str := r.RemoteAddr

	fmt.Printf("Here 1: %v\n", auth)

	if len(auth) > bearerLen+1 && auth[:bearerLen] == bearer {
		t, err := jwt.Parse(auth[bearerLen+1:], func(token *jwt.Token) (interface{}, error) {
			// Always check the signing method
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				//fmt.Printf("Here 3: %v\n", "not good")
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// Return the key for validation
			fmt.Printf("Here username: %v\n", token.Claims["username"])
			fmt.Printf("Here token: %v\n", token)

			username := token.Claims["username"].(string)

			secretkey, err := mykeys(username)

			if err == nil {
				return []byte(secretkey), nil
			} else {
				return []byte(secretkey), nil
			}
		})

		if err == nil && t.Valid {
			outMap["status"] = true
			outMap["username"] = t.Claims["username"].(string)
			outMap["msg"] = "Good"
		} else {
			outMap["msg"] = "Bad"
		}
	} else {
	}

	outJSON, err := json.Marshal(outMap)

	if err != nil {
		w.Write([]byte("{\"status\": false}"))
	}

	w.Write([]byte(outJSON))
}

func serveJWT(w http.ResponseWriter, r *http.Request) {
	outMap := map[string]interface{}{
		"status": false,
	}

	username := "username1"
	secretkey, err := mykeys(username)

	if err != nil {
		w.Write([]byte("{\"status\": false}"))
	}

	token := genJsonWebToken(username, secretkey)

	outMap["token"] = token
	outMap["status"] = true

	outJSON, err := json.Marshal(outMap)

	if err != nil {
		w.Write([]byte("{\"status\": false}"))
	}

	w.Write([]byte(outJSON))
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

	http.HandleFunc("/validatejwt", serveValidateJWT)
	http.HandleFunc("/jwt", serveJWT)

	http.HandleFunc("/", serveHome)

	err := http.ListenAndServe(":80", nil)

	if err != nil {
		fmt.Printf("main(): %s\n", err)
		log.Fatal("ListenAndServe: ", err)
	}
}
