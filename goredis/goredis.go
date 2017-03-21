package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/go-redis/redis"
)

type Person struct {
	Name string
	Age  int
}

var (
	client *redis.ClusterClient
)

func main() {
	var err error

	//client := redis.NewClient(&redis.Options{
	//	Addr:     "localhost:30001", // Redis default port is 6379
	//	Password: "",                // no password set
	//	DB:       0,                 // use default DB
	//})

	client = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{":30001", ":30002", ":30003", ":30004", ":30005", ":30006"},
	})

	pong, err := client.Ping().Result()
	fmt.Printf("Pong: %#v and %#v\n", pong, err)

	DemoString()
	DemoStruct()
}

func DemoString() {
	if str, err := client.Get("hello").Result(); err != nil {
		fmt.Printf("err: %s\n", err)
	} else {
		fmt.Printf("hello: %s\n", str)
	}
}

func DemoStruct() {
	p1 := Person{Name: "Jun", Age: 19}
	buf := new(bytes.Buffer)
	gob.NewEncoder(buf).Encode(&p1)

	if err := client.HSet("test1", "person", buf.Bytes()).Err(); err != nil {
		fmt.Printf("HSet err: %s\n", err)
		return
	}

	if b, err := client.HGet("test1", "person").Bytes(); err != nil {
		fmt.Printf("HGet err: %s\n", err)
		return
	} else {
		p2 := Person{}
		gob.NewDecoder(bytes.NewReader(b)).Decode(&p2)
		fmt.Printf("Person: %#v\n", p2)
	}
}
