package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

type Employee struct {
	Status    bool
	Firstname string
	LastName  string
	Age       int
}

func main() {
	// this example shows how to loop through the field names of the struct and set field value by field name.
	emp := Employee{}

	empFields := reflect.TypeOf(Employee{})

	for i := 0; i < empFields.NumField(); i++ {
		fieldName := empFields.Field(i).Name
		fieldType := empFields.Field(i).Type

		fmt.Printf("%s\t%s\n", fieldName, fieldType)

		if fieldName == "Status" || fieldName == "Age" {
			continue
		}

		fieldVal := "Test " + strconv.Itoa(i)

		reflect.ValueOf(&emp).Elem().FieldByName(fieldName).SetString(fieldVal)
	}

	empJSON, err := json.MarshalIndent(emp, "", " ")

	if err == nil {
		fmt.Printf("%s\n", empJSON)
	}
}
