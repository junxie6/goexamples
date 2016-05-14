package main

import (
	"encoding/json"
	"fmt"
)

// convJSONToMap shows how to decode the nested JSON string into a map by doing type assertion.
// Without type assertion, you will get invalid operation: (type interface {} does not support indexing)
func convJSONToMap() {
	strJSON := `{"100":{"DealerName":"Company A","DealerCountry":"USA","ItemArr":[{"LineNum":0,"IdItem":10,"ItemNum":"A0"},{"LineNum":1,"IdItem":11,"ItemNum":"A1"}]},"101":{"DealerName":"Company B","DealerCountry":"Canada","ItemArr":[{"LineNum":0,"IdItem":10,"ItemNum":"A0"},{"LineNum":1,"IdItem":11,"ItemNum":"A1"}]}}`

	var data map[string]interface{}

	if err := json.Unmarshal([]byte(strJSON), &data); err != nil {
	} else {
		for IdDealer, v0 := range data {
			DealerName := v0.(map[string]interface{})["DealerName"]
			DealerCountry := v0.(map[string]interface{})["DealerCountry"]

			fmt.Printf("%v\t%v\t%v\n", IdDealer, DealerName, DealerCountry)

			ItemArr := v0.(map[string]interface{})["ItemArr"]

			for k1, v1 := range ItemArr.([]interface{}) {
				LineNum := v1.(map[string]interface{})["LineNum"]
				IdItem := v1.(map[string]interface{})["IdItem"]
				ItemNum := v1.(map[string]interface{})["ItemNum"]
				fmt.Printf("\tLine[%v]:\t%v\t%v\t%v\n", k1, LineNum, IdItem, ItemNum)
			}

			// If the underlying type is unknown, a type switch determines the type:
			switch v := ItemArr.(type) {
			case int:
				fmt.Printf("It's int: %v\n", v)
			case float64:
				fmt.Printf("It's float64: %v\n", v)
			case string:
				fmt.Printf("It's string: %v\n", v)
			case []interface{}:
				fmt.Printf("It's []interface{}: %v\n", v)
			default:
				fmt.Printf("It is not one of the types above: %v\n", v)
			}
		}
	}
}

func main() {
	convJSONToMap()
}
