package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

// convJSONToMap shows how to decode the nested JSON string into a map by doing type assertion.
// Without type assertion, you will get invalid operation: (type interface {} does not support indexing)
func convJSONToMap() {
	strJSON := `{"100":{"DealerName":"Company A","DealerCountry":"USA","ItemArr":[{"LineNum":0,"IdItem":10,"ItemNum":"A0"},{"LineNum":1,"IdItem":11,"ItemNum":"A1"}]},"101":{"DealerName":"Company B","DealerCountry":"Canada","ItemArr":[{"LineNum":0,"IdItem":10,"ItemNum":"A0"},{"LineNum":1,"IdItem":11,"ItemNum":"A1"}]}}`

	var data map[string]interface{}

	if err := json.Unmarshal([]byte(strJSON), &data); err != nil {
	} else {
		for IdDealer, v0 := range data {
			//
			DealerName := v0.(map[string]interface{})["DealerName"]
			DealerCountry := v0.(map[string]interface{})["DealerCountry"]

			fmt.Printf("\n%v\t%v\t%v\n", IdDealer, DealerName, DealerCountry)

			ItemArr := v0.(map[string]interface{})["ItemArr"]

			for k1, v1 := range ItemArr.([]interface{}) {
				LineNum := v1.(map[string]interface{})["LineNum"]
				IdItem := v1.(map[string]interface{})["IdItem"]
				ItemNum := v1.(map[string]interface{})["ItemNum"]
				fmt.Printf("\tLine[%v]:\t%v\t%v\t%v\n", k1, LineNum, IdItem, ItemNum)
			}

			// Optional: Use reflect to determine the type:
			fmt.Printf("Type: %v\n", reflect.TypeOf(ItemArr))

			// Optional: If the underlying type is unknown, use a type switch to determine the type:
			switch v := ItemArr.(type) {
			case int:
				fmt.Printf("Type: int: %v\n", v)
			case float64:
				fmt.Printf("Type: float64: %v\n", v)
			case string:
				fmt.Printf("Type: string: %v\n", v)
			case []interface{}:
				fmt.Printf("Type: []interface{}: %v\n", v)
			default:
				fmt.Printf("Type: It is not one of the types above: %v\n", v)
			}
		}
	}
}

func convJSONTOStruct() {
	type SOLine struct {
		LineNum int
		ItemNum string
		Price   float64
	}
	var data = struct {
		Status    bool
		PONum     string
		SOLineArr []SOLine
		IDOrder   int
		IDDealer  int
		Test      string
		test2     string
	}{}

	jsonStr := `{"Status":true,"IDOrder":5,"IDDealer":99,"PONum":"My PO","SOLineArr":[{"LineNum":0,"ItemNum":"AAA0","Price":99},{"LineNum":1,"ItemNum":"AAA1","Price":99}]}`

	json.NewDecoder(bytes.NewBufferString(jsonStr)).Decode(&data)

	fmt.Printf("IDOrder: %v\n", data.IDOrder)
	fmt.Printf("IDDealer: %v\n", data.IDDealer)

	for _, v := range data.SOLineArr {
		fmt.Printf("SOLineArr.ItemNum: %v\n", v.ItemNum)
	}
}

func main() {
	//convJSONToMap()
	convJSONTOStruct()
}
