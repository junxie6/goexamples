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

// example1 - Simple
func example1() {
	var postURL = "http://127.0.0.1/"

	//
	formData := url.Values{"action": {"createOrder"}}
	//formData := make(url.Values)

	formData.Set("idDealer", "916") // Set sets the key to value. It replaces any existing values.
	formData.Add("poNum", "test")   // Add adds the value to key. It appends to any existing values associated with key.

	//
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.PostForm(postURL, formData)

	if err != nil {
		log.Printf("%v", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Printf("%v", err)
	}

	if resp.StatusCode != http.StatusOK {
		// do something
	}

	fmt.Printf("%v\t%s", resp.Status, body)
}

// example2 - Advanced
func example2() {
	var postURL = "http://127.0.0.1/"

	formData := url.Values{"action": {"createOrder"}}
	//formData := make(url.Values)

	formData.Set("idDealer", "916") // Set sets the key to value. It replaces any existing values.
	formData.Add("poNum", "test")   // Add adds the value to key. It appends to any existing values associated with key.

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

	//
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("%v", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Printf("%v", err)
	}

	if resp.StatusCode != http.StatusOK {
		// do something
	}

	fmt.Printf("%v\t%v\n", resp.Status, resp.Header)
	fmt.Printf("%s\n", body)
}

// example3 - Advanced - post JSON
func example3() {
	postURL := "http://127.0.0.1:8080/so?act=SaveSO"

	var jsonStr = `{"Status":true,"Data":{"so":{"IDOrder":1,"Status":1,"Created":"123","Changed":"456"}}}`

	// TODO: Is strings.NewReader() faster than bytes.NewBufferString()?
	req, err := http.NewRequest("POST", postURL, bytes.NewBuffer([]byte(jsonStr)))
	//req, err := http.NewRequest("POST", postURL, strings.NewReader(jsonStr))

	if err != nil {
		log.Printf("%v", err)
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	//
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("%v", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Printf("%v", err)
	}

	if resp.StatusCode != http.StatusOK {
		// do something
	}

	fmt.Printf("%v\t%v\n", resp.Status, resp.Header)
	fmt.Printf("%s\n", body)
}

func main() {
	//example1()
	//example2()
	example3()
}
