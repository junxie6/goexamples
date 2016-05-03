package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"log"
)

var (
	db *sql.DB
)

func initDB() {
	var err error
	db, err = sql.Open("mysql", "MyUser:MyPassword@tcp(localhost:3306)/MyDB")

	if err != nil {
		log.Fatalf("Error on initializing database connection: %v", err.Error())
	}

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()

	if err != nil {
		log.Fatalf("Error on opening database connection: %v", err.Error())
	}
}

func testInsert(username string, plaintextPassword string) {
	var hashPassword string

	if s, err := HashPassword(plaintextPassword); err != nil {
		log.Printf("%v", err)
	} else {
		hashPassword = string(s)
	}

	stmt, err := db.Prepare("INSERT INTO users (username, password) VALUES (?, ?)")

	if err != nil {
		log.Printf("Error: %v", err)
	}

	defer stmt.Close()

	res, err := stmt.Exec(username, hashPassword)

	if err != nil {
		log.Printf("Error: %v", err)
	}

	lastID, err := res.LastInsertId()

	if err != nil {
		log.Printf("Error: %v", err)
	}

	rowCnt, err := res.RowsAffected()

	if err != nil {
		log.Printf("Error: %v", err)
	}

	fmt.Printf("testInsert: ID = %d, affected = %d\n", lastID, rowCnt)
}

func testSelectMultipleRowsV1() {
	// query
	var (
		uid      int
		username string
	)

	stmt, err := db.Prepare("SELECT uid, username FROM users WHERE ?")

	if err != nil {
		log.Printf("Error: %v", err)
	}

	defer stmt.Close()

	rows, err := stmt.Query(1)

	if err != nil {
		log.Printf("Error: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&uid, &username)

		if err != nil {
			log.Printf("Error: %v", err)
		}

		fmt.Printf("testSelectMultipleRowsV1: %d, %s\n", uid, username)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error: %v", err)
	}
}

func testSelectMultipleRowsV2() {
	// query
	var (
		uid      int
		username string
	)

	// Under the hood, db.Query() actually prepares, executes, and closes a prepared statement.
	// That's three round-trips to the database. If you're not careful, you can triple the number of database interactions your application makes!
	rows, err := db.Query("SELECT uid, username FROM users WHERE ?", 1)

	if err != nil {
		log.Printf("Error: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&uid, &username)

		if err != nil {
			log.Printf("Error: %v", err)
		}

		fmt.Printf("testSelectMultipleRowsV2: %d, %s\n", uid, username)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error: %v", err)
	}
}

func testSelectSingleRow() {
	var (
		uid  int
		name string
	)

	err := db.QueryRow("SELECT uid, username FROM users WHERE username =  ?", "jun").Scan(&uid, &name)

	if err != nil {
		log.Printf("Error: %v", err)
	}

	fmt.Printf("testSelectSingleRow: %d, %s\n", uid, name)
}

// HashPassword ...
func HashPassword(plaintextPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(plaintextPassword), bcrypt.DefaultCost)
}

// ValidatePassword ...
func ValidatePassword(hashed string, plaintextPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plaintextPassword))
}

func main() {
	initDB()
	defer db.Close()

	testInsert("test1", "test1")
	testSelectMultipleRowsV1()
	testSelectMultipleRowsV2()
	testSelectSingleRow()
}
