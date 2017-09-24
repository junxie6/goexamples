package model

import (
	"database/sql"
	//"fmt"
	"time"
)

import (
	"github.com/jmoiron/sqlx"
)

type Ticket struct {
	IDTicket  uint      `db:"IDTicket"`
	IDProject uint      `db:"IDProject"`
	IDUser    uint      `db:"IDUser"`
	Subject   string    `db:"Subject"`
	Changed   time.Time `db:"Changed"`
}

func ListTicket() ([]Ticket, error) {
	var err error
	ticketArr := make([]Ticket, 0)

	sq := "SELECT "
	sq += "  IDTicket "
	sq += ", IDProject "
	sq += ", IDUser "
	sq += ", Subject "
	sq += ", Changed "
	sq += "FROM ticket "

	err = db.Select(&ticketArr, sq)

	if err == sql.ErrNoRows {
		return ticketArr, nil
	} else if err != nil {
		return nil, err
	}

	return ticketArr, nil
}

func GenerateTicket() {
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

	//
	var stmt *sqlx.Stmt

	sq := "INSERT INTO ticket (IDProject, IDUser, Subject, Changed) "
	sq += "VALUES (?, ?, ?, ?) "

	stmt, err = db.Preparex(sq)

	if err != nil {
		fmt.Printf("Err: %s\n", err.Error())
		return
	}

	stmt.Exec()
}

//func GetTicket() {
//	var err error
//	var rows *sqlx.Rows
//
//	ticket := Ticket{}
//
//	sq := "SELECT "
//	sq += "  IDTicket "
//	sq += ", IDProject "
//	sq += ", IDUser "
//	sq += ", Subject "
//	sq += ", Changed "
//	sq += "FROM ticket "
//
//	rows, err = db.Queryx(sq)
//
//	if err != nil {
//		fmt.Printf("Err: %s\n", err.Error())
//		return
//	}
//
//	for rows.Next() {
//		err := rows.StructScan(&ticket)
//
//		if err != nil {
//			fmt.Printf("Err: %s\n", err.Error())
//			return
//		}
//
//		fmt.Printf("%#v\n", ticket)
//	}
//}
