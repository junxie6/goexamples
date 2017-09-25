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
	var IDTicket uint

	sq := "INSERT INTO ticket_keyword (IDTicket, IDKeyword) "
	sq += "VALUES (?, ?) "

	for c := 1; c < 1000000; {
		tx := db.MustBegin()

		stmt, err = tx.Preparex(sq)

		if err != nil {
			fmt.Printf("Err: %s\n", err.Error())
			return
		}

		i := 0

		for ; i < 1000; i++ {
			IDTicket = uint(c) + uint(i)
			IDKeywordMap := make(map[int]bool)

			for kc := 0; kc < 5; kc++ {
				randomNum = util.RandomNumber(0, 30)
				IDKeywordMap[randomNum] = true
			}

			for IDKeyword, _ := range IDKeywordMap {
				_, err = stmt.Exec(IDTicket, IDKeyword)

				if err != nil {
					fmt.Printf("Err: %s\n", err.Error())
					return
				}
			}

			//fmt.Printf("%d ", IDTicket)
		}

		c += i

		tx.Commit()
	}
}
