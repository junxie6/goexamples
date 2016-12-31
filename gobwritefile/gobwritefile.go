package main

import (
	"encoding/gob"
	"fmt"
	"os"
)

// Reference:
// https://www.socketloop.com/tutorials/golang-saving-and-reading-file-with-gob
func main() {
	filename := "integerdata.gob"

	WriteGob(filename)
	ReadGob(filename)
}

type Person struct {
	Name string
	Age  int
}

// WriteGob ...
func WriteGob(filename string) {
	data := []Person{
		Person{Name: "AAA", Age: 16},
		Person{Name: "BBB", Age: 18},
	}

	// create a file
	dataFile, err := os.Create(filename)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer dataFile.Close()

	// serialize the data
	dataEncoder := gob.NewEncoder(dataFile)

	if err := dataEncoder.Encode(data); err != nil {
		fmt.Printf("err: %v\n", err)
	}
}

// ReadGob ...
func ReadGob(filename string) {
	var data []Person

	// open data file
	dataFile, err := os.Open(filename)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer dataFile.Close()

	dataDecoder := gob.NewDecoder(dataFile)

	if err := dataDecoder.Decode(&data); err != nil {
		fmt.Printf("err: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("data: %#v\n", data)
}
