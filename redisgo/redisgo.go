package main

import (
	"github.com/garyburd/redigo/redis"
	"log"
	"net/http"
)

var Conn redis.Conn

func srvSet(w http.ResponseWriter, r *http.Request) {
	// set
	Conn.Do("SET", "message1", "Hello World")
}

func srvGet(w http.ResponseWriter, r *http.Request) {

	// get
	if world, err := redis.String(Conn.Do("GET", "message1")); err != nil {
		w.Write([]byte("key not found"))
	} else {
		w.Write([]byte(world))
	}
}

func Dial() (redis.Conn, error) {
	return redis.Dial("tcp", ":6379")
}

func main() {
	var err error

	Conn, err = Dial()

	if err != nil {
		log.Printf("%v", err.Error())
	}

	defer Conn.Close()

	http.HandleFunc("/Set", srvSet)
	http.HandleFunc("/Get", srvGet)

	http.ListenAndServe(":8080", nil)
}
