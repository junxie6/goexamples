package model

import (
//"database/sql"
//"fmt"
)

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	db *sqlx.DB
)

func Open() error {
	var err error
	db, err = sqlx.Connect("mysql", "exp:exp@tcp(localhost:3306)/exp?parseTime=true")

	if err != nil {
		return err
	}

	return nil
}

func Close() error {
	return db.Close()
}
