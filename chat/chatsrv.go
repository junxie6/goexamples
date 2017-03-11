// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 254.
//!+

// Chat is a server that lets clients chat with each other.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
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
	clients := make(map[client]bool) // all connected clients

	for {
		select {
		case cli := <-entering:
			for c := range clients {
				cli.ch <- c.name + " is in room."
			}

			clients[cli] = true
		case cli := <-leaving:
			delete(clients, cli)
			close(cli.ch)
		case msg := <-messaging:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				cli.ch <- msg
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
		messaging <- cli.name + ": " + input.Text()
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

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")

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
