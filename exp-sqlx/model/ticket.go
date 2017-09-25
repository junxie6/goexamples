package model

import (
	"database/sql"
	"fmt"
	"time"
)

import (
	"github.com/jmoiron/sqlx"
	"github.com/junhsieh/util"
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
	var projectArr []Project

	projectArr, err = ListProject()

	if err != nil {
		fmt.Printf("Err: %s\n", err.Error())
		return
	}

	// Insert new record
	var stmt *sqlx.Stmt
	var IDUser uint
	var IDProject uint
	var randomNum int
	var Subject string
	numOfProject := len(projectArr)

	sq := "INSERT INTO ticket (IDProject, IDUser, Subject, Changed) "
	sq += "VALUES (?, ?, ?, NOW()) "

	stmt, err = db.Preparex(sq)

	if err != nil {
		fmt.Printf("Err: %s\n", err.Error())
		return
	}

	for i := 0; i < 1000000; i++ {
		randomNum = util.RandomNumber(0, numOfProject)
		IDProject = projectArr[randomNum].IDProject
		IDUser = uint(util.RandomNumber(1, 10))
		Subject = fmt.Sprintf("%s did not work due to incorrect handle the function in case of error need to fix it", projectArr[randomNum].Name)

		_, err = stmt.Exec(IDProject, IDUser, Subject)

		if err != nil {
			fmt.Printf("Err: %s\n", err.Error())
			return
		}
	}
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
