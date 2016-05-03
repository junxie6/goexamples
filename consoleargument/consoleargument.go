package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatal("filename not specified")
	}

	for k, v := range flag.Args() {
		fmt.Printf("%d: %s\n", k, v)
	}

	filename := flag.Args()[0]
	fmt.Printf("Filename: %s\n", filename)
}
