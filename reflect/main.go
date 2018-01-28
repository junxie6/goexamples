package main

import (
	"fmt"
	"reflect"
)

const tagName = "validate"

type User struct {
	Id    int    `validate:"-"`
	Name  string `validate:"presence,min=2,max=32"`
	Email string `validate:"email,required"`
}

func main() {
	user := User{
		Id:    1,
		Name:  "John Doe",
		Email: "john@example",
	}

	ShowFieldNameTypeValueTag(&user)

	//rv := reflect.TypeOf(&user)

	//fmt.Printf("HERE: %#v\n", rv.Kind().String())

	////rv := reflect.ValueOf(&user)
	//for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
	//	//	fmt.Println(rv.Kind(), rv.Type())
	//	rv = rv.Elem()
	//}

	//t := rv
	//t.Elem()

	//t := reflect.TypeOf(user)

	//fmt.Println("Type:", t.Name())
	//fmt.Println("Kind:", t.Kind())

	// Iterate over all available fields and read the tag value
	//for i := 0; i < t.NumField(); i++ {
	//	// Get the field, returns https://golang.org/pkg/reflect/#StructField
	//	valueField := t.Field(i)
	//	typeField := t.Type().Field(i)
	//	tag := typeField.Tag.Get(tagName)

	//	// Get the field tag value
	//	//field := t.Field(i)
	//	//tag := field.Tag.Get(tagName)
	//	//tag := field.Tag.Get(tagName)

	//	//fmt.Printf("%d. %v (%v): , tag: '%v'\n", i+1, field.Name, field.Type.Name(), tag)
	//	fmt.Printf("%d. %v (%v): %#v, tag: '%v'\n", i+1, typeField.Name, valueField.Type(), valueField.Interface(), tag)
	//	//field.Name
	//	//fmt.Printf("%#v\n", field.Interface())
	//}

	//rv := reflect.TypeOf(&user)
	//fmt.Printf("%#v %#v\n", rv.Kind().String(), reflect.Ptr)

	//return

	//rv := reflect.ValueOf(&user)
	//for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
	//	//	fmt.Println(rv.Kind(), rv.Type())
	//	rv = rv.Elem()
	//}

	//fmt.Println(rv.Kind())
	//fmt.Println(rv.Kind(), rv.Type())

	//t2 := rv.Type()

	//for i := 0; i < t2.NumField(); i++ {
	//	field := t.Field(i)
	//	fmt.Printf("%#v\n", field.Name)
	//}
}

func ShowFieldNameTypeValueTag(v interface{}) {
	t := reflect.ValueOf(v).Elem()

	fmt.Printf("%#v ============\n", t.Kind().String())

	for i := 0; i < t.NumField(); i++ {
		valueField := t.Field(i)
		typeField := t.Type().Field(i)
		tag := typeField.Tag.Get(tagName)

		fmt.Printf("%d. %v (%v): %#v, tag: '%v'\n", i+1, typeField.Name, valueField.Type(), valueField.Interface(), tag)
	}
}
