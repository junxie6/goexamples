package main

import (
	"fmt"
	"reflect"
)

const tagName = "validate"

type User struct {
	ID            int          `validate:"-"`
	Name          string       `validate:"Name,min=2,max=32"`
	Email         string       `validate:"Email,required"`
	Bag           Bag          `jjj:"Bag"`
	CreditCardArr []CreditCard `validate:"CreditCardArr"`
}

type Bag struct {
	Name   string `jjj:"Bag"`
	Weight uint   `jjj:"Weight"`
}

type CreditCard struct {
	Number string `validate:"Number"`
}

func main() {
	user := User{
		ID:    1,
		Name:  "Jun Xie",
		Email: "jun@example.com",
	}

	ShowFieldNameTypeValueTagV1(&user)
}

func ShowFieldNameTypeValueTagV1(v interface{}) {
	t := reflect.ValueOf(v).Elem()
	typeOfT := t.Type()

	fmt.Printf("%#v %#v============\n", t.Kind().String(), typeOfT.Name())

	for i := 0; i < t.NumField(); i++ {
		valueField := t.Field(i)
		typeField := typeOfT.Field(i)
		tag := typeField.Tag.Get(tagName)

		fmt.Printf("%d. %v (%v|%v): %#v, tag: %v\n", i+1, typeField.Name, valueField.Kind().String(), valueField.Type(), valueField.Interface(), tag)

		if valueField.Kind() == reflect.Struct {
			fieldPtr := valueField.Addr()
			ShowFieldNameTypeValueTagV1(fieldPtr.Interface())
		}
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
