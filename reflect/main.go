package main

// Reference:
// https://stackoverflow.com/questions/34493062/how-to-reflect-struct-recursive-in-golang
// https://stackoverflow.com/questions/20496585/reflecting-the-fields-of-a-emptys-slices-underyling-type
// https://gist.github.com/hvoecking/10772475

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
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
	Table         [][]string
}

type Bag struct {
	Name   string `jjj:"Bag"`
	Weight uint   `jjj:"Weight"`
}

type CreditCard struct {
	Number      string   `validate:"Number"`
	ExpireDate  string   `jjj:"ExpireDate"`
	BranchIDArr []string `jjj:"BranchIDArr"`
	MonthArr    []int
	AccountArr  []Account
	Person      Person
}

type Person struct {
	Name        string
	Age         uint
	LanguageArr []string
	CarArr      []Car
}

type Car struct {
	Model     string
	YearBuilt uint
}

type Account struct {
	Num     string
	Balance float64
}

func main() {
	user := User{
		ID:    1,
		Name:  "Jun Xie",
		Email: "jun@example.com",
		Bag: Bag{
			Name:   "LV",
			Weight: 79,
		},
		Table: [][]string{
			[]string{"R0 C0", "R0 C1"},
			[]string{"R1 C0", "R1 C1"},
		},
		CreditCardArr: []CreditCard{
			CreditCard{
				Number:      "123",
				ExpireDate:  "2018-01-29",
				BranchIDArr: []string{"111", "222", "333"},
				MonthArr:    []int{1, 3, 5},
				AccountArr: []Account{
					Account{
						Num:     "111111111111111",
						Balance: 5999,
					},
					Account{
						Num:     "222222222222222",
						Balance: 6999,
					},
				},
				Person: Person{
					Name:        "jun 0",
					Age:         12,
					LanguageArr: []string{"English", "Chinese"},
					CarArr: []Car{
						Car{
							Model:     "X5",
							YearBuilt: 1999,
						},
						Car{
							Model:     "X6",
							YearBuilt: 1998,
						},
					},
				},
			},
			CreditCard{
				Number:      "456",
				ExpireDate:  "2018-01-29",
				BranchIDArr: []string{"444", "555", "666"},
				MonthArr:    []int{2, 4, 6},
				AccountArr: []Account{
					Account{
						Num:     "33333333333333",
						Balance: 5999,
					},
					Account{
						Num:     "444444444444444",
						Balance: 6999,
					},
				},
				Person: Person{
					Name:        "jun 1",
					Age:         13,
					LanguageArr: []string{"English", "Chinese"},
					CarArr: []Car{
						Car{
							Model:     "X5",
							YearBuilt: 1999,
						},
						Car{
							Model:     "X6",
							YearBuilt: 1998,
						},
					},
				},
			},
		},
	}

	//translated := translate(&user)
	//fmt.Println("original:  ", original, "->", (*original.Payload), "->", (*original.Payload).(B).Ptr)
	//fmt.Println("translated:", translated, "->", (*translated.(D).Payload), "->", (*(translated.(D).Payload)).(B).Ptr)
	//fmt.Printf("%#v\n", translated)

	//ShowFieldNameTypeValueTagV1(&user)

	test := make(map[string]interface{})
	Flatten(&user, test, "", false)
	//fmt.Printf("%#v\n", test)
	ObjectToJSON(test, true)
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

		if valueField.Kind() == reflect.Slice || valueField.Kind() == reflect.Array {
			for ii := 0; ii < valueField.Len(); ii++ {
				//fmt.Printf("HHHH: %#v\n", valueField.Index(ii).Interface())
				fieldPtr := valueField.Index(ii).Addr()
				ShowFieldNameTypeValueTagV1(fieldPtr.Interface())
			}
		}
	}
}

func Flatten(v interface{}, data map[string]interface{}, parentStr string, isParentASliceOrStruct bool) {
	t := reflect.ValueOf(v).Elem()
	typeOfT := t.Type()

	//fmt.Printf("%#v %#v============\n", t.Kind().String(), typeOfT.Name())

	asdf2 := typeOfT.Name()

	for i := 0; i < t.NumField(); i++ {
		valueField := t.Field(i)
		typeField := typeOfT.Field(i)
		//tag := typeField.Tag.Get(tagName)

		// FieldName, FieldSysType, FieldUserType, FieldValue, Tag
		//fmt.Printf("%d. %v (%v|%v): %#v, tag: %v\n", i+1, typeField.Name, valueField.Kind().String(), valueField.Type(), valueField.Interface(), tag)

		//fmt.Printf("%v.%v\n", typeOfT.Name(), typeField.Name)

		var key string

		if isParentASliceOrStruct == false {
			key = parentStr + asdf2 + "." + typeField.Name
		} else {
			key = parentStr + "." + typeField.Name
		}

		if valueField.Kind() == reflect.Struct {
			fieldPtr := valueField.Addr()
			Flatten(fieldPtr.Interface(), data, key, true)
			continue
		}

		if valueField.Kind() == reflect.Slice || valueField.Kind() == reflect.Array {
			objKey := key // User.CreditCardArr

			for ii := 0; ii < valueField.Len(); ii++ {
				key = objKey + "[" + strconv.Itoa(ii) + "]"

				if valueField.Index(ii).Kind() == reflect.Struct {
					fieldPtr := valueField.Index(ii).Addr()
					Flatten(fieldPtr.Interface(), data, key, true)
					continue
				}

				if valueField.Index(ii).Kind() == reflect.Slice || valueField.Index(ii).Kind() == reflect.Array {
					objKey2 := key
					valueField2 := valueField.Index(ii)

					for iii := 0; iii < valueField2.Len(); iii++ {
						key = objKey2 + "[" + strconv.Itoa(iii) + "]"

						if valueField2.Index(iii).Kind() == reflect.Struct {
							fieldPtr2 := valueField2.Index(iii).Addr()
							Flatten(fieldPtr2.Interface(), data, key, true)
							continue
						}

						data[key] = valueField2.Index(iii).Interface()
					}
					continue
				}

				// For string, or int ??
				data[key] = valueField.Index(ii).Interface()

				//data[key] = ""

				//fmt.Printf("HHHH: %#v\n", valueField.Index(ii).Interface())
			}

			continue
		}

		data[key] = valueField.Interface()
		//switch vvv := valueField.Interface().(type) {
		//case int:
		//	data[key] = vvv
		//case string:
		//	data[key] = vvv
		//case float64:
		//	data[key] = vvv
		//}
	}
}

func ObjectToJSON(v interface{}, isIndent bool) {
	var byteArr []byte
	var err error

	if isIndent == true {
		byteArr, err = json.MarshalIndent(v, "", "    ")
	} else {
		byteArr, err = json.Marshal(v)
	}

	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	fmt.Printf("v: %s\n", string(byteArr))
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
