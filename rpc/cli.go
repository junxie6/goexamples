package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"time"
)

func main() {
	// Connection without setting a timeout:
	conn, err := rpc.Dial("tcp", "localhost:8080")

	if err != nil {
		log.Fatal(err)
	}

	// If you want to add timeout:
	//conn0, err := net.DialTimeout("tcp", "localhost:8080", 5*time.Second)

	//if err != nil {
	//	log.Fatal("dialing:", err)
	//}

	//conn := rpc.NewClient(conn0)

	//
	in := bufio.NewReader(os.Stdin)

	for {
		line, _, err := in.ReadLine()

		if err != nil {
			log.Printf("ReadLine: %#v\n", err)
		}

		reply := new([]byte)

		err = conn.Call("Myfunc.Echo", line, reply)

		// If an error is returned, the reply parameter will not be sent back to the client.
		if err != nil {
			log.Printf("Call: %#v\n", err)
		}

		fmt.Printf("Replied: %#v\n", string(*reply))
	}

	conn.Close()
}
