package main

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
)

func main() {
	client, err := rpc.Dial("tcp", "localhost:8080")

	if err != nil {
		log.Fatal(err)
	}

	in := bufio.NewReader(os.Stdin)

	for {
		line, _, err := in.ReadLine()

		if err != nil {
			log.Printf("ReadLine: %#v\n", err)
		}

		reply := new([]byte)

		err = client.Call("Myfunc.Echo", line, reply)

		// If an error is returned, the reply parameter will not be sent back to the client.
		if err != nil {
			log.Printf("Call: %#v\n", err)
		}

		fmt.Printf("Replied: %#v\n", string(*reply))
	}
}
