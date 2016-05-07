package main

import (
	"encoding/json"
	"fmt"
)

func convJSONToMap() {
	strJSON := `{"100":{"DealerName":"Company A","DealerCountry":"USA","ItemArr":[{"LineNum":0,"IdItem":10,"ItemNum":"A0"},{"LineNum":1,"IdItem":11,"ItemNum":"A1"}]},"101":{"DealerName":"Company B","DealerCountry":"Canada","ItemArr":[{"LineNum":0,"IdItem":10,"ItemNum":"A0"},{"LineNum":1,"IdItem":11,"ItemNum":"A1"}]}}`

	var data map[string]interface{}

	_ = json.Unmarshal([]byte(strJSON), &data)

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
	}
}

func main() {
	convJSONToMap()
}
