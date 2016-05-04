package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func exampleConvMapToJSONOutputBrowser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	orderMap := map[string]interface{}{
		"status":  true,
		"name":    "bot",
		"age":     18,
		"itemArr": []interface{}{},
	}

	orderMap["email"] = "bot@example.com"

	item0 := map[string]interface{}{
		"lineNum":   0,
		"itemNum":   "A0",
		"itemPrice": 90,
	}
	item1 := map[string]interface{}{
		"lineNum":      1,
		"itemNum":      "A1",
		"itemPrice":    91,
		"accessoryArr": []interface{}{},
	}

	item1["accessoryArr"] = append(item1["accessoryArr"].([]interface{}), map[string]interface{}{"name": "bag", "price": 10})
	item1["accessoryArr"] = append(item1["accessoryArr"].([]interface{}), map[string]interface{}{"name": "tray", "price": 12})

	item2 := map[string]interface{}{
		"lineNum":   2,
		"itemNum":   "A2",
		"itemPrice": 92,
	}

	orderMap["itemArr"] = append(orderMap["itemArr"].([]interface{}), item0)
	orderMap["itemArr"] = append(orderMap["itemArr"].([]interface{}), item1)
	orderMap["itemArr"] = append(orderMap["itemArr"].([]interface{}), item2)

	orderJson, err := json.Marshal(orderMap)

	if err == nil {
		w.Write([]byte(orderJson))
	} else {
		w.Write([]byte(`{"status": false}`))
	}
}

func main() {
	http.HandleFunc("/exampleConvMapToJSONOutputBrowser", exampleConvMapToJSONOutputBrowser)

	if err := http.ListenAndServe(":80", nil); err != nil {
		fmt.Printf("main(): %s\n", err)
		log.Fatal("ListenAndServe: ", err)
	}
}
