package main

import (
	"fmt"
)

import (
	"github.com/junhsieh/goexamples/gobgzipv2/gobgzipv2"
)

// Person ...
type Person struct {
	Name   string
	Age    int
	Money  float64
	Houses []House
}

// House ...
type House struct {
	StreetNum     int
	StreetName    string
	PropertyValue float64
}

func main() {
	person := GetPerson()
	person2 := Person{}

	if b, err := gobgzipv2.EncodeGobThenGzip(&person); err != nil {
		fmt.Printf("err: %v\n\n", err)
	} else if err := gobgzipv2.UngzipThenDecodeGob(b, &person2); err != nil {
		fmt.Printf("err: %v\n\n", err)
	} else {
		fmt.Printf("gob encoded: %v\n\n", b)
		fmt.Printf("person2: %#v\n\n", person2)
	}
}

func GetPerson() Person {
	return Person{
		Name:  "AAA",
		Age:   123,
		Money: 123456.789,
		Houses: []House{
			House{
				StreetNum:     111,
				StreetName:    "Hello 111",
				PropertyValue: 123456.789,
			},
			House{
				StreetNum:     222,
				StreetName:    "Hello 222",
				PropertyValue: 987654.321,
			},
		},
	}
}
