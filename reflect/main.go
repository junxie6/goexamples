package main

// Reference:
// https://stackoverflow.com/questions/34493062/how-to-reflect-struct-recursive-in-golang
// https://stackoverflow.com/questions/20496585/reflecting-the-fields-of-a-emptys-slices-underyling-type
// https://gist.github.com/hvoecking/10772475

import (
	"fmt"
	"reflect"
)

const tagName = "validate"

var dict = map[string]string{
	"Hello!":                 "Hallo!",
	"What's up?":             "Was geht?",
	"translate this":         "übersetze dies",
	"point here":             "zeige hier her",
	"translate this as well": "übersetze dies auch...",
	"and one more":           "und noch eins",
	"deep":                   "tief",
}

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

	translated := translate(&user)
	//fmt.Println("original:  ", original, "->", (*original.Payload), "->", (*original.Payload).(B).Ptr)
	//fmt.Println("translated:", translated, "->", (*translated.(D).Payload), "->", (*(translated.(D).Payload)).(B).Ptr)
	fmt.Printf("%#v\n", translated)
	//ShowFieldNameTypeValueTagV1(&user)
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

func translate(obj interface{}) interface{} {
	// Wrap the original in a reflect.Value
	original := reflect.ValueOf(obj)

	copy := reflect.New(original.Type()).Elem()
	translateRecursive(copy, original)

	// Remove the reflection wrapper
	return copy.Interface()
}

func translateRecursive(copy, original reflect.Value) {
	switch original.Kind() {
	// The first cases handle nested structures and translate them recursively

	// If it is a pointer we need to unwrap and call once again
	case reflect.Ptr:
		// To get the actual value of the original we have to call Elem()
		// At the same time this unwraps the pointer so we don't end up in
		// an infinite recursion
		originalValue := original.Elem()
		// Check if the pointer is nil
		if !originalValue.IsValid() {
			return
		}
		// Allocate a new object and set the pointer to it
		copy.Set(reflect.New(originalValue.Type()))
		// Unwrap the newly created pointer
		translateRecursive(copy.Elem(), originalValue)

	// If it is an interface (which is very similar to a pointer), do basically the
	// same as for the pointer. Though a pointer is not the same as an interface so
	// note that we have to call Elem() after creating a new object because otherwise
	// we would end up with an actual pointer
	case reflect.Interface:
		// Get rid of the wrapping interface
		originalValue := original.Elem()
		// Create a new object. Now new gives us a pointer, but we want the value it
		// points to, so we have to call Elem() to unwrap it
		copyValue := reflect.New(originalValue.Type()).Elem()
		translateRecursive(copyValue, originalValue)
		copy.Set(copyValue)

	// If it is a struct we translate each field
	case reflect.Struct:
		for i := 0; i < original.NumField(); i += 1 {
			translateRecursive(copy.Field(i), original.Field(i))
		}

	// If it is a slice we create a new slice and translate each element
	case reflect.Slice:
		copy.Set(reflect.MakeSlice(original.Type(), original.Len(), original.Cap()))
		for i := 0; i < original.Len(); i += 1 {
			translateRecursive(copy.Index(i), original.Index(i))
		}

	// If it is a map we create a new map and translate each value
	case reflect.Map:
		copy.Set(reflect.MakeMap(original.Type()))
		for _, key := range original.MapKeys() {
			originalValue := original.MapIndex(key)
			// New gives us a pointer, but again we want the value
			copyValue := reflect.New(originalValue.Type()).Elem()
			translateRecursive(copyValue, originalValue)
			copy.SetMapIndex(key, copyValue)
		}

	// Otherwise we cannot traverse anywhere so this finishes the the recursion

	// If it is a string translate it (yay finally we're doing what we came for)
	case reflect.String:
		translatedString := dict[original.Interface().(string)]
		copy.SetString(translatedString)

	// And everything else will simply be taken from the original
	default:
		copy.Set(original)
	}
}
