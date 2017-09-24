// User > Role > Permission > Project
package main

import (
	//"encoding/json"
	"fmt"
)

import (
	"github.com/junhsieh/goexamples/exp-sqlx/model"
)

func GetTicketList() {
	var err error
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

func GetProjectList() {
	var err error
	var projectArr []model.Project

	projectArr, err = model.ListProject()

	if err != nil {
		fmt.Printf("Err: %s\n", err.Error())
		return
	}

	for _, project := range projectArr {
		fmt.Printf("%#v\n", project)
	}
}

func main() {
	var err error

	if err = model.Open(); err != nil {
		fmt.Printf("Err: %s\n", err.Error())
		return
	}

	defer model.Close()

	//
	//GetTicketList()

	// Generate tickets
	model.GenerateTicket()
}
