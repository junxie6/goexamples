package main

import (
	//"database/sql"
	"fmt"
)

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	db *sqlx.DB
)

type Ticket struct {
	IDTicket  uint   `db:"IDTicket"`
	IDProject uint   `db:"IDProject"`
	IDUser    uint   `db:"IDUser"`
	Subject   string `db:"Subject"`
}

func GetTicket() {
	var err error
	var rows *sqlx.Rows

	ticket := Ticket{}

	sq := "SELECT "
	sq += "  IDTicket "
	sq += ", IDProject "
	sq += ", IDUser "
	sq += ", Subject "
	sq += "FROM ticket "

	rows, err = db.Queryx(sq)

	if err != nil {
		fmt.Printf("Err: %s\n", err.Error())
		return
	}

	for rows.Next() {
		err := rows.StructScan(&ticket)

		if err != nil {
			fmt.Printf("Err: %s\n", err.Error())
			return
		}

		fmt.Printf("%#v\n", ticket)
	}
}

func main() {
	var err error

	db, err = sqlx.Connect("mysql", "exp:exp@tcp(localhost:3306)/exp")

	if err != nil {
		fmt.Printf("Err: %s\n", err.Error())
		return
	}

	defer db.Close()

	//
	GetTicket()
}
