package main

import (
	//"bytes"
	//"encoding/gob"
	"fmt"
	"golang.org/x/net/publicsuffix"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
)

func main() {
	// To store the cookie:
	// https://github.com/juju/persistent-cookiejar

	cookie := example_httppost()
	example_httpget(cookie)
}

func example_httppost() *cookiejar.Jar {
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}

	jar, err := cookiejar.New(&options)

	if err != nil {
		log.Fatal(err)
	}

	// Declare http client
	client := &http.Client{
		Jar: jar,
	}

	// Declare post data
	formData := make(url.Values)
	//formData := url.Values{}
	//PostData := strings.NewReader("act=set&name=bot")

	formData.Set("act", "set")
	formData.Set("name", "bot")

	// Declare HTTP Method and Url
	var postURL = "http://erp.local:8080/"

	// TODO: Is strings.NewReader() faster than bytes.NewBufferString()?
	req, err := http.NewRequest("POST", postURL, strings.NewReader(formData.Encode()))
	//req, err := http.NewRequest("POST", postURL, bytes.NewBufferString(formData.Encode()))

	if err != nil {
		log.Printf("http.NewRequest: %v", err.Error())
	}

	// Must set Content-Type. Otherwise, web server will not pick up the data.
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Add("Content-Length", strconv.Itoa(len(formData.Encode())))

	// Set cookie
	//req.Header.Set("Cookie", "name=MySession; count=1")

	resp, err := client.Do(req)

	// Read response
	if data, err := ioutil.ReadAll(resp.Body); err != nil {
		fmt.Printf("error = %s \n", err)
	} else {
		// Print response
		log.Printf("Response = %s", string(data))

		log.Printf("Cookie: %v", jar)
	}

	return jar
}

func example_httpget(jar *cookiejar.Jar) {
	// Declare http client
	client := &http.Client{
		Jar: jar,
	}

	// Declare post data
	//formData := make(url.Values)

	// Declare HTTP Method and Url
	var postURL = "http://erp.local:8080/"

	req, err := http.NewRequest("GET", postURL+"?act=Test2", strings.NewReader("act=Test"))

	if err != nil {
		log.Printf("http.NewRequest: %v", err.Error())
	}

	// Set cookie
	//req.Header.Set("Cookie", "name=MySession; count=1")

	resp, err := client.Do(req)

	// Read response
	if data, err := ioutil.ReadAll(resp.Body); err != nil {
		fmt.Printf("error = %s \n", err)
	} else {
		// Print response
		log.Printf("Response2 = %s", string(data))
		log.Printf("Cookie2: %v", jar)
	}
}
