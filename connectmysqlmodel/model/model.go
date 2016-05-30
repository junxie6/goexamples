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

// txWrapper usage example:
//so := SO{}
//
//ifArr := []interface{}{
//	&so,
//}
//
//err := txWrapper(func(tx *sql.Tx, ifArr []interface{}) error {
//	if err := ifArr[0].(*SO).lockSO(tx); err != nil {
//		return err
//	}
//
//	return nil
//}, ifArr)

func txWrapper(txFunc func(*sql.Tx, []interface{}) error, ifArr []interface{}) (retErr error) {
	tx, retErr := db.Begin()

	if retErr != nil {
		return
	}

	// Although tx could be accessed inside the defer function
	// it is good practice to pass it through argument. Because:
	// "deferred function's arguments are evaluated when the defer statement is evaluated."
	// "…they may refer to variables defined in a surrounding function.
	// Those variables are then shared between the surrounding function and the function literal, and they survive as long as they are accessible."
	// http://www.blevesearch.com/news/Deferred-Cleanup,-Checking-Errors,-and-Potential-Problems/
	// https://blog.golang.org/defer-panic-and-recover
	defer func(tx *sql.Tx) {
		// If the current goroutine is panicking, a call to recover will capture the value given to panic and resume normal execution.
		// We capcture the error before making the final decision (rollback/commit).
		if err := recover(); err != nil {
			switch v := err.(type) {
			case error:
				retErr = v
			case string:
				retErr = errors.New(v)
			default:
				retErr = fmt.Errorf("%v", v)
			}
		}

		if retErr != nil {
			// Do not try to assign the return value of tx.Rollback to err directily.
			// If you do, you would invert the meaning of the final return value.
			// Because if Rollback succeeds, it would return nil.
			if err := tx.Rollback(); err != nil {
				retErr = err
			}
			return
		}

		retErr = tx.Commit()
	}(tx)

	return txFunc(tx, ifArr) // the return value of txFunc will be set to retErr.
}
