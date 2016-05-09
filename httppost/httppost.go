package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	//"strings"
	"bytes"
	"strconv"
	"time"
)

func example1() {
	var postURL = "http://127.0.0.1/"
	client := &http.Client{Timeout: 20 * time.Second}

	formData := url.Values{"action": {"createOrder"}}
	formData.Set("idDealer", "916") // Set sets the key to value. It replaces any existing values.
	formData.Add("poNum", "test")   // Add adds the value to key. It appends to any existing values associated with key.

	resp, err := client.PostForm(postURL, formData)

	if err != nil {
		log.Printf("%v", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Printf("%v", err)
	}

	fmt.Printf("%v\t%s", resp.Status, body)
}

func example2() {
	var postURL = "http://127.0.0.1/"

	formData := url.Values{"action": {"createOrder"}}
	formData.Set("idDealer", "916") // Set sets the key to value. It replaces any existing values.
	formData.Add("poNum", "test")   // Add adds the value to key. It appends to any existing values associated with key.

	client := &http.Client{}

	// first way
	req, err := http.NewRequest("POST", postURL, bytes.NewBufferString(formData.Encode()))

	// second way
	//req, err := http.NewRequest("POST", postURL, strings.NewReader(formData.Encode()))

	if err != nil {
		log.Printf("%v", err)
	}

	req.Header.Add("Authorization", "Bearer: ASDF")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(formData.Encode())))

	resp, _ := client.Do(req)

	if err != nil {
		log.Printf("%v", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Printf("%v", err)
	}

	fmt.Printf("%v\t%s", resp.Status, body)
}

func main() {
	//example1()
	example2()
}
