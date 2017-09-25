package model

import (
	//"database/sql"
	"fmt"
	//"time"
)

import (
	"github.com/jmoiron/sqlx"
	"github.com/junhsieh/util"
)

type TicketKeyword struct {
	IDTicket  uint `db:"IDTicket"`
	IDKeyword uint `db:"IDKeyword"`
}

func GenerateTicketKeyword() {
	var err error

	// Insert new record
	var stmt *sqlx.Stmt
	var randomNum int

	sq := "INSERT INTO ticket_keyword (IDTicket, IDKeyword) "
	sq += "VALUES (?, ?) "

	tx := db.MustBegin()

	stmt, err = tx.Preparex(sq)

	if err != nil {
		fmt.Printf("Err: %s\n", err.Error())
		return
	}

	for i := 0; i < 10; i++ {
		IDKeywordMap := make(map[int]bool)

		for kc := 0; kc < 5; kc++ {
			randomNum = util.RandomNumber(0, 60)
			IDKeywordMap[randomNum] = true
		}

		for IDKeyword, _ := range IDKeywordMap {
			_, err = stmt.Exec(i, IDKeyword)

			if err != nil {
				fmt.Printf("Err: %s\n", err.Error())
				return
			}
		}

	}

	tx.Commit()
}
