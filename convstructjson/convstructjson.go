package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Todo struct {
	Name      string `json:"name"` // By adding struct tags you can control exactly what and how your struct field name will be marshalled to JSON.
	Completed bool
	due       time.Time
}

type Todos []Todo

type Item struct {
	ItemName  string
	ItemPrice int
}

type MyOrder struct {
	Status       bool
	CustomerName string
	ItemArr      []Item
}

// This example convert the struct to a JSON string.
func exampleConvStructToJSONOutputScreen1() {
	todos := Todos{
		Todo{Name: "Write presentation", Due: time.Now()},
		Todo{Name: "Host meetup", Due: time.Now()},
	}

	if data, err := json.Marshal(todos); err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	} else {
		fmt.Printf("json.Marshal:\n%s\n\n", data)
	}

	// produces neatly indented output
	if data, err := json.MarshalIndent(todos, "", " "); err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	} else {
		fmt.Printf("json.MarshalIndent:\n%s\n\n", data)
	}
}

// This example convert the struct to a JSON string and output to screen directly.
func exampleConvStructToJSONOutputScreen2() {
	todos := Todos{
		Todo{Name: "Write presentation", Due: time.Now()},
		Todo{Name: "Host meetup", Due: time.Now()},
	}

	fmt.Println("json.NewEncoder:")
	if err := json.NewEncoder(os.Stdout).Encode(todos); err != nil {
		fmt.Println(err)
	}
}

func getMyOrder() MyOrder {
	item0 := Item{
		ItemName:  "Computer",
		ItemPrice: 90,
	}
	item2 := Item{ItemName: "Motherboard", ItemPrice: 92}
	item3 := Item{ItemName: "Hard drive", ItemPrice: 93}
	item4 := Item{ItemName: "Mouse", ItemPrice: 94}
	item5 := Item{ItemName: "Monitor", ItemPrice: 95}
	item6 := Item{ItemName: "Book", ItemPrice: 96}
	item7 := Item{ItemName: "Desk", ItemPrice: 97}
	item8 := Item{ItemName: "Door", ItemPrice: 98}
	item9 := Item{ItemName: "Car", ItemPrice: 99}

	order1 := MyOrder{
		Status:  true,
		ItemArr: []Item{item0, Item{ItemName: "Pen", ItemPrice: 91}}, // initialize with two items.
	}

	order1.CustomerName = "Jun"

	// first way of appending the items to the ItemArr slice.
	order1.ItemArr = append(order1.ItemArr, item2, item3)

	// second way of appending the items to the ItemArr slice.
	// The three dots is called variadic parameter / Ellipsis meaning unpack this slice into a set of variadic arguments.
	order1.ItemArr = append(order1.ItemArr, []Item{item4, item5}...)

	// third way of appending the items to the ItemArr slice.
	itemArr := []Item{item6, item7}
	order1.ItemArr = append(order1.ItemArr, itemArr...)

	// fourth way of appending the items to the ItemArr slice.
	order1.ItemArr = append(order1.ItemArr, item8)
	order1.ItemArr = append(order1.ItemArr, item9)

	return order1
}

func exampleConvStructToJSONOutputBrowser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	order1 := getMyOrder()
	encoder := json.NewEncoder(w)

	if err := encoder.Encode(order1); err != nil {
		fmt.Println(err)
	}
}

func main() {
	exampleConvStructToJSONOutputScreen1() // output to screen
	exampleConvStructToJSONOutputScreen2() // output to screen

	http.HandleFunc("/exampleConvStructToJSONOutputBrowser", exampleConvStructToJSONOutputBrowser)

	if err := http.ListenAndServe(":80", nil); err != nil {
		fmt.Printf("main(): %s\n", err)
		log.Fatal("ListenAndServe: ", err)
	}
}
