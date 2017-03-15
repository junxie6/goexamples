package main

import (
	//"errors"
	"fmt"
	"log"
	"net"
	"net/rpc"
	//"time"
)

type Myfunc int

func (m *Myfunc) Echo(line []byte, ack *[]byte) error {
	//time.Sleep(5 * time.Second)

	*ack = []byte("GOOD")
	fmt.Printf("Your entered: %#v _ %#v\n", string(line), string(*ack))

	//return errors.New("test err")
	return nil
}

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:8080")

	if err != nil {
		log.Fatal(err)
	}

	rpc.Register(new(Myfunc))

	rpc.Accept(listener)

	// If you want to execute some other behaviour in-between the request being accepted and it being served)
	// you can change the code as follows to swap out rpc.Accept(listener) above.
	//for {
	//	conn, err := listener.Accept()
	//	if err != nil {
	//		// handle error
	//	}

	//	go rpc.ServeConn(conn)
	//}
}
