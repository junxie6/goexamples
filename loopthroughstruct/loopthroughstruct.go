package main

import (
	"fmt"
	"reflect"
)

type person struct {
	Status    bool
	Firstname string
	Lastname  string
	Age       int
}

func main() {
	typ := reflect.TypeOf(person{})

	fmt.Printf("%v\n", typ)

	for i := 0; i < typ.NumField(); i++ {
		fmt.Printf("%v %v\n", typ.Field(i).Name, typ.Field(i).Type)
	}
}
