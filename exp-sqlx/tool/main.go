package main

import (
	"fmt"
)

import (
	"github.com/junhsieh/goexamples/exp-sqlx/model"
)

// User > Role > Permission > Project
func main() {
	var err error

	if err = model.Open(); err != nil {
		fmt.Printf("Err: %s\n", err.Error())
		return
	}

	defer model.Close()

	//
	var ticketArr []model.Ticket

	ticketArr, err = model.ListTicket()

	if err != nil {
		fmt.Printf("Err: %s\n", err.Error())
		return
	}

	for _, ticket := range ticketArr {
		fmt.Printf("%#v\n", ticket)
	}
}
