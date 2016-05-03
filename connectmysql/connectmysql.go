package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var (
	db *sql.DB
)

func test() {
	// query
	var (
		id   int
		name string
	)

	rows, err := db.Query("SELECT uid, username from users where uid = ?", 1)

	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	//defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &name)

		if err != nil {
			log.Fatalf("Error: %s", err)
		}

		fmt.Println(id, name)
	}

	err = rows.Err()

	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}

func initDB() {
	var err error
	db, err = sql.Open("mysql", "MyUser:MyPassword@tcp(localhost:3306)/MyDB")

	if err != nil {
		log.Fatalf("Error on initializing database connection: %s", err.Error())
	}

	//defer db.Close()

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()

	if err != nil {
		log.Fatalf("Error on opening database connection: %s", err.Error())
	}
}

func main() {
	initDB()
	test()
}
