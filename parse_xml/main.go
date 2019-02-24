package main

// Reference:
// https://tutorialedge.net/golang/parsing-xml-with-golang/
// https://golang.org/pkg/encoding/xml/#example_Unmarshal

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

// our struct which contains the complete
// array of all Users in the file
type Users struct {
	XMLName xml.Name `xml:"users"`
	Users   []User   `xml:"user"`
}

// the user struct, this contains our
// Type attribute, our user's name and
// a social struct which will contain all
// our social links
type User struct {
	XMLName   xml.Name `xml:"user"`
	Type      string   `xml:"type,attr"`
	AB        string   `xml:"b,attr"`
	Languages []string `xml:"languages>name"`
	Name      string   `xml:"name"`
	Social    Social   `xml:"social"`
}

// a simple struct which contains all our
// social links
type Social struct {
	XMLName  xml.Name `xml:"social"`
	Facebook string   `xml:"facebook"`
	Twitter  string   `xml:"twitter"`
	Youtube  string   `xml:"youtube"`
}

func main() {
	var xmlFile *os.File
	var err error

	// Open our xmlFile
	if xmlFile, err = os.Open("demo.xml"); err != nil {
		fmt.Println(err)
		return
	}

	// defer the closing of our xmlFile so that we can parse it later on
	defer xmlFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(xmlFile)

	// we initialize our Users array
	var users Users

	// we unmarshal our byteArray which contains our
	// xmlFiles content into 'users' which we defined above
	if err = xml.Unmarshal(byteValue, &users); err != nil {
		fmt.Println(err)
	}

	// we iterate through every user within our users array and
	// print out the user Type, their name, and their facebook url
	// as just an example
	for i := 0; i < len(users.Users); i++ {
		fmt.Println("User Type: " + users.Users[i].Type)
		fmt.Println("AB: " + users.Users[i].AB)
		fmt.Println("User Name: " + users.Users[i].Name)
		fmt.Printf("Languages: %v\n", users.Users[i].Languages)
		fmt.Println("Facebook Url: " + users.Users[i].Social.Facebook)
	}
}
