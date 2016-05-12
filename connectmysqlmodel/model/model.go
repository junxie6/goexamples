package model

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB
)

func InitDB(dsn string) (*sql.DB, error) {
	var err error

	// The database/sql package manages the connection pooling automatically for you.
	// sql.Open(..) returns a handle which represents a connection pool, not a single connection.
	// The database/sql package automatically opens a new connection if all connections in the pool are busy.
	// Reference: http://stackoverflow.com/questions/17376207/how-to-share-mysql-connection-between-http-goroutines
	//db, err = sql.Open("mysql", dsn)
	//db, err = sql.Open("mysql", "MyUser:MyPassword@tcp(localhost:3306)/MyDB?tx_isolation='READ-COMMITTED'") // optional

	if db, err = sql.Open("mysql", dsn); err != nil {
		// Error on initializing database connection.
		return nil, err
	} else if err := db.Ping(); err != nil {
		// Error on opening database connection.
		return nil, err
	}

	return db, nil
}

// txEnd handles the transaction Rollback and Commit logic by checking errors.
func txEnd(tx *sql.Tx, errArrPtr *[]error) {
	if len(*errArrPtr) > 0 {
		if err := tx.Rollback(); err != nil {
			*errArrPtr = append(*errArrPtr, err)
		}
	} else if err := tx.Commit(); err != nil {
		*errArrPtr = append(*errArrPtr, err)
	}
}
