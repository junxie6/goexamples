package main

import (
	"fmt"
	"github.com/go-redis/redis"
)

func main() {
	var err error

	//client := redis.NewClient(&redis.Options{
	//	Addr:     "localhost:30001", // Redis default port is 6379
	//	Password: "",                // no password set
	//	DB:       0,                 // use default DB
	//})

	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{":30001", ":30002", ":30003", ":30004", ":30005", ":30006"},
	})

	pong, err := client.Ping().Result()
	fmt.Printf("Pong: %#v and %#v\n", pong, err)

	var str string

	str, err = client.Get("hello").Result()

	if err != nil {
		fmt.Printf("err: %s\n", err)
	}

	fmt.Printf("hello: %s\n", str)

	str, err = client.Get("foo").Result()

	if err != nil {
		fmt.Printf("err: %s\n", err)
	}

	fmt.Printf("foo: %s\n", str)
}
