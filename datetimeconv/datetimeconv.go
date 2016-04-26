package main

import (
	"flag"
	"fmt"
	"gopkg.in/redis.v3"
	"log"
	"reflect"
	"strconv"
	"time"
)

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		log.Fatal("Please specify Redis server IP")
	}

	redisServer := flag.Args()[0]

	client := redis.NewClient(&redis.Options{
		Addr:     redisServer + ":6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	t := time.Now()

	// returns t as a Unix time, the number of nanoseconds elapsed since January 1, 1970 UTC
	t1 := t.UnixNano()

	// returns the string representation of t1 in 16 base
	t2 := strconv.FormatInt(t1, 16)

	t3, _ := strconv.ParseInt(t2, 16, 64)

	t4 := time.Unix(0, t3)

	fmt.Printf("t1: %v, %v\n", t1, reflect.TypeOf(t1))
	fmt.Printf("t2: %v, %v\n", t2, reflect.TypeOf(t2))
	fmt.Printf("t3: %v, %v\n", t3, reflect.TypeOf(t3))
	fmt.Printf("t4: %v, %v\n", t4, reflect.TypeOf(t4))

	err := client.Set("MyLastMod", t1, 0).Err()

	if err != nil {
		fmt.Printf("ERR: %s\n", err)
	}

	s, err := client.Get("MyLastMod").Int64()

	client.Close()

	if err != nil {
		fmt.Printf("ERR: %s\n", err)
	}

	fmt.Printf("Last Modified Time: %v %v\n", s, reflect.TypeOf(s))
}
