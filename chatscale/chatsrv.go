// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 254.
//!+

// Chat is a server that lets clients chat with each other.
package main

import (
	"bufio"
	//"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

type client struct {
	name string
	ip   string
	ch   chan<- string
}

var (
	entering  = make(chan client)
	leaving   = make(chan client)
	messaging = make(chan string)
)

func broadcaster() {
	clients := make(map[string]client) // all connected clients

	for {
		select {
		case cli := <-entering:
			for _, c := range clients {
				cli.ch <- c.name + " is in room."
			}

			clients[cli.name] = cli
		case cli := <-leaving:
			delete(clients, cli.name)
			close(cli.ch)
		case msg := <-messaging:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			//for _, cli := range clients {
			//	cli.ch <- msg
			//}
			strArr := strings.Split(msg, ":")
			// 0: sender
			// 1: receiver
			// 2: msg

			if *addr == ":8080" {
				fmt.Printf("db1: %#v\n", strArr)
				strArr = []string{
					"bb",
					"aa",
					"muhaha",
				}
				if cli, ok := clients[strArr[1]]; ok {
					cli.ch <- strArr[2]
				}
			}

			if *addr == ":8081" && len(strArr) >= 3 {
				fmt.Printf("db2: %#v\n", strArr)
				go talkToSrv(":8080", strArr[1], strArr[2])
			}
		}
	}
}

func handleConn(conn net.Conn) {
	cli := client{}
	ch := make(chan string) // outgoing client messages

	go clientWriter(conn, ch)

	cli.ip = conn.RemoteAddr().String()
	cli.ch = ch

	cli.ch <- "You are " + cli.ip

	cli.ch <- "Please enter your name:"
	input2 := bufio.NewScanner(conn)

	for input2.Scan() {
		cli.name = input2.Text()
		break
	}

	//
	messaging <- cli.name + "@" + cli.ip + " has arrived."
	entering <- cli

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messaging <- cli.name + ":" + input.Text()
	}
	// NOTE: ignoring potential errors from input.Err()

	leaving <- cli
	messaging <- cli.name + "@" + cli.ip + " has left."

	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintf(conn, "%s\n", msg) // NOTE: ignoring network errors
	}
}

func talkToSrv(addr string, name string, msg string) {
	conn, err := net.Dial("tcp", "localhost"+addr)
	if err != nil {
		log.Fatal(err)
	}

	// for incoming message
	//done := make(chan struct{})
	//go func() {
	//	io.Copy(os.Stdout, conn) // NOTE: ignoring errors
	//	done <- struct{}{}       // signal the main goroutine
	//}()

	// for outgoing message
	//mustCopy(conn, bytes.NewBufferString("srv2"))
	mustCopy(conn, strings.NewReader("srv2"))
	mustCopy(conn, strings.NewReader(name+":"+msg))

	conn.Close()
	//<-done // wait for background goroutine to finish
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}

var (
	addr *string = flag.String("addr", ":8080", "The addr of the application.")
)

func main() {
	// flags
	//var addr = flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse()

	//
	listener, err := net.Listen("tcp", "localhost"+*addr)

	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Printf("Err: %v", err)
			continue
		}

		go handleConn(conn)
	}
}
