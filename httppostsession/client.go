package main

import (
	"fmt"
	"golang.org/x/net/publicsuffix"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	//"bytes"
)

func main() {
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
	//PostData := strings.NewReader("act=set&name=bot")
	postData := url.Values{}
	postData.Set("act", "set")
	postData.Set("name", "bot")

	// Declare HTTP Method and Url
	var postURL = "http://erp.local:8080/"

	req, err := http.NewRequest("POST", postURL, strings.NewReader(postData.Encode()))
	//req, err := http.NewRequest("POST", postURL, bytes.NewBufferString(postData.Encode()))

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(postData.Encode())))

	// Set cookie
	//req.Header.Set("Cookie", "name=MySession; count=1")

	resp, err := client.Do(req)

	// Read response
	data, err := ioutil.ReadAll(resp.Body)

	// error handle
	if err != nil {
		fmt.Printf("error = %s \n", err)
	}

	// Print response
	fmt.Printf("Response = %s", string(data))
	fmt.Printf("Cookie: %v", jar)
}
