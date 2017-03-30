package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"strings"
)

var csvFile = strings.NewReader(`Name,Age
jun,19
jun,20
jun,21
`)

func main() {
	r := csv.NewReader(csvFile)
	rows, err := r.ReadAll()
	if err != nil {
		log.Fatalf("error reading all lines: %v", err)
	}

	for k, row := range rows {
		fmt.Printf("%d: %v\n", k, row)
	}
}
