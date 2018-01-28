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
		Name:  "Jun Xie",
		Email: "jun@example.com",
	}

	ShowFieldNameTypeValueTagV1(&user)
}

func ShowFieldNameTypeValueTagV1(v interface{}) {
	t := reflect.ValueOf(v).Elem()

	fmt.Printf("%#v ============\n", t.Kind().String())

	typeOfT := t.Type()

	for i := 0; i < t.NumField(); i++ {
		valueField := t.Field(i)
		typeField := typeOfT.Field(i)
		tag := typeField.Tag.Get(tagName)

		fmt.Printf("%d. %v (%v): %#v, tag: '%v'\n", i+1, typeField.Name, valueField.Type(), valueField.Interface(), tag)
	}
}

func ShowFieldNameTypeValueTagV2(v interface{}) {
	t := reflect.TypeOf(v)

	for t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface {
		t = t.Elem()
	}

	fmt.Printf("%#v ============\n", t.Kind().String())

	for i := 0; i < t.NumField(); i++ {
		valueField := t.Field(i)
		typeField := valueField.Type

		typeName0 := typeField.Name()
		//typeName1 := typeField.Kind().String()
		//typeName2 := typeField.String()

		tag := valueField.Tag.Get(tagName)

		fmt.Printf("%d. %v (%v), tag: '%v'\n", i+1, valueField.Name, typeName0, tag)
	}
}

func ShowFieldNameTypeValueTagV3(v interface{}) {
	//rv := reflect.ValueOf(&user)
}

func ShowFieldNameTypeValueTagV4(v interface{}) {
	//reflect.Indirect()
	//fmt.Println("Indirect type is:", reflect.Indirect(reflect.ValueOf(v)).Elem().Type()) // prints main.CustomStruct
	//fmt.Println("Indirect value type is:", reflect.Indirect(reflect.ValueOf(v)).Elem().Kind()) // prints struct
}
