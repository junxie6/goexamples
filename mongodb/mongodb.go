package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type Person struct {
	Name string
	Age  int
}

func main() {
	conn, err := mgo.Dial("127.0.0.1:27017")

	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Optional. Switch the session to a monotonic behavior.
	conn.SetMode(mgo.Monotonic, true)

	session := conn.DB("mydb").C("mycol")

	// insert
	people := []interface{}{
		Person{"jun0", 29},
		Person{"jun1", 29},
	}

	//if err := session.Insert(people...); err != nil {
	if err := session.Insert(people...); err != nil {
		log.Fatal(err)
	}

	// select one record
	person1 := Person{}

	if err := session.Find(bson.M{"name": "jun0"}).One(&person1); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Single: %#v\n", person1)

	// select multiple records
	personArr := []Person{}

	// The option 'i' of bson.RegEx is for case insensitive matching.
	if err := session.Find(bson.M{"name": bson.RegEx{"^JUN", "i"}}).All(&personArr); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Multiple: %#v\n", personArr)
}
