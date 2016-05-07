package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"log"
	"reflect"
	"strings"
)

var (
	db *sql.DB
)

func initDB() {
	var err error

	// The database/sql package manages the connection pooling automatically for you.
	// sql.Open(..) returns a handle which represents a connection pool, not a single connection.
	// The database/sql package automatically opens a new connection if all connections in the pool are busy.
	// Reference: http://stackoverflow.com/questions/17376207/how-to-share-mysql-connection-between-http-goroutines
	db, err = sql.Open("mysql", "MyUser:MyPassword@tcp(localhost:3306)/MyDB")
	//db, err = sql.Open("mysql", "MyUser:MyPassword@tcp(localhost:3306)/MyDB?tx_isolation='READ-COMMITTED'") // optional

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

func testSelectMultipleRowsV3(optArr map[string]interface{}) {
	// queries
	query := []string{}
	param := []interface{}{}

	if val, ok := optArr["idOrder"]; ok {
		query = append(query, "salesOrder.idOrder >= ?")
		param = append(param, val)
	}

	// The first character of the field name must be in upper case. Otherwise, you would get:
	// panic: reflect.Value.Interface: cannot return value obtained from unexported field or method
	var sqlField = struct {
		IdOrder int
		Uid     int
		Changed string
	}{}

	var rowArr []interface{}

	sqlFieldArrPtr := StrutToSliceOfFieldAddress(&sqlField)

	sql := ""
	sql += "SELECT "
	sql += "  salesOrder.idOrder "
	sql += ", salesOrder.uid "
	sql += ", salesOrder.changed "
	sql += "FROM salesOrder "
	sql += "WHERE " + strings.Join(query, " AND ") + " "
	sql += "ORDER BY salesOrder.idOrder "

	stmt, err := db.Prepare(sql)
	if err != nil {
		log.Printf("Error: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(param...)

	if err != nil {
		log.Printf("Error: %v", err)
	}

	defer rows.Close()

	if err != nil {
		log.Printf("Error: %v", err)
	}

	//sqlFields, err := rows.Columns()

	for rows.Next() {
		err := rows.Scan(sqlFieldArrPtr...)

		if err != nil {
			log.Printf("Error: %v", err)
		}

		// Show the type of each struct field
		f1 := reflect.TypeOf(sqlField.IdOrder)
		f2 := reflect.TypeOf(sqlField.Uid)
		f3 := reflect.TypeOf(sqlField.Changed)
		fmt.Printf("Type: %v\t%v\t%v\n", f1, f2, f3)

		// Show the value of each field
		fmt.Printf("Row: %v\t%v\t%v\n\n", sqlField.IdOrder, sqlField.Uid, sqlField.Changed)

		rowArr = append(rowArr, sqlField)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error: %v", err)
	}

	// produces neatly indented output
	if data, err := json.MarshalIndent(rowArr, "", " "); err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	} else {
		fmt.Printf("json.MarshalIndent:\n%s\n\n", data)
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

func StrutToSliceOfFieldAddress(s interface{}) []interface{} {
	fieldArr := reflect.ValueOf(s).Elem()

	fieldAddrArr := make([]interface{}, fieldArr.NumField())

	for i := 0; i < fieldArr.NumField(); i++ {
		f := fieldArr.Field(i)
		fieldAddrArr[i] = f.Addr().Interface()
	}

	return fieldAddrArr
}

func main() {
	initDB()
	defer db.Close()

	// this example shows how to insert data.
	//testInsert("test1", "test1")

	// this example shows how to select multiple rows
	//testSelectMultipleRowsV1()

	// this example shows how to select multiple rows
	//testSelectMultipleRowsV2()

	// this example shows how to dynamically assign a list of field name to the rows.Scan() function.
	optArr := map[string]interface{}{}
	optArr["idOrder"] = 1
	testSelectMultipleRowsV3(optArr)

	// this example shows how to select single row.
	testSelectSingleRow()
}
