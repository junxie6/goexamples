package model

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/junhsieh/goexamples/util"
	"log"
	"strings"
	"time"
)

const (
	// SOStatusInit is when a dummy SO created.
	SOStatusInit = 0

	// SOStatusEdit is when user is inputting SO data.
	SOStatusEdit = 1

	// SOStatusConfirmed is when SO is confirmed to be processed.
	SOStatusConfirmed = 2

	// SOStatusPartialReversed is when some SP of the SO are not reversed (these unreversed SP are to be invoiced).
	SOStatusPartialReversed = 3

	// SOStatusFullReversed is when all SP of the SO are reversed.
	SOStatusFullReversed = 4

	// SOStatusFulfilled is when all ordered items are fulfilled.
	SOStatusFulfilled = 5
)

// SO is SalesOrder
type SO struct {
	IDOrder   int
	Status    int
	Created   string
	Changed   string
	SOInfo    SOInfo
	SOAddr    SOAddr
	SOLineArr []SOLine
	PackArr   []SP
}

// SOInfo ...
type SOInfo struct {
	IDOrder  int
	IDDealer int
	PONum    string
}

// SOAddr ...
type SOAddr struct {
	IDOrder     int
	IDShipAddr  int
	DealerName  string
	BillContact string
	BillStreet1 string
	BillStreet2 string
	BillCity    string
	BillState   string
	BillZip     string
	BillCountry string
	ShipContact string
	ShipStreet1 string
	ShipStreet2 string
	ShipCity    string
	ShipState   string
	ShipZip     string
	ShipCountry string
}

// SOLine ...
type SOLine struct {
	IDOrder      int
	LineNum      int
	IDWarehouse  int
	IDItem       int
	OrderedQty   int
	ShippedQty   int
	BackOrderQty int
	Changed      string
}

func (so *SO) Insert(tx *sql.Tx, errArrPtr *[]error) {
	so.Status = 1

	stmt, err := tx.Prepare("INSERT INTO SO (Status, changed) VALUES (?, NOW())")

	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		defer stmt.Close()

		if rs, err := stmt.Exec(so.Status); err != nil {
			log.Printf("Error: %v", err)
		} else if lastID, err := rs.LastInsertId(); err != nil {
			log.Printf("Error: %v", err)
		} else {
			log.Printf("lastID: %v", lastID)
		}
	}
}

// editSO runs the queries in transaction.
func (so *SO) EditSO() []error {
	errArr := []error{}

	if tx, err := db.Begin(); err != nil {
		errArr = append(errArr, err)
	} else {
		// Lock SO - considering whether to add "locked IDOrder rows" checking as well.
		if so.lockSO(tx, &errArr); len(errArr) == 0 {
			// [DEBUG] Sleep
			if so.SOInfo.PONum == "lock" {
				time.Sleep(15 * time.Second)
			}

			// more actions once it's locked
			//so.editSORaw(tx, &errArr)
			//so.Insert(tx, &errArr)
		}

		txEnd(tx, &errArr)
	}

	return errArr
}

func (so *SO) lockSO(tx *sql.Tx, errArrPtr *[]error) {
	if lockSORaw := so.lockSORaw(tx, errArrPtr); len(lockSORaw) == 0 {
		*errArrPtr = append(*errArrPtr, fmt.Errorf("IDOrderArr is empty. No row is locked"))
	}
}

// lockSORaw ...
func (so *SO) lockSORaw(tx *sql.Tx, errArrPtr *[]error) []int {
	optArr := map[string]interface{}{}

	optArr["StatusArr"] = []int{SOStatusInit, SOStatusEdit, SOStatusConfirmed}

	outArr := []int{}
	query := []string{}
	param := []interface{}{}

	if val, ok := optArr["StatusArr"]; ok {
		query = append(query, "SO.Status IN ("+util.Placeholder(len(val.([]int)))+")")

		for _, v := range val.([]int) {
			param = append(param, v)
		}
	}

	if len(query) == 0 {
		*errArrPtr = append(*errArrPtr, fmt.Errorf("SQL query has no WHERE condition"))
	}

	sql := "SELECT "
	sql += "  SO.IDOrder "
	sql += "FROM SO "
	sql += "WHERE " + strings.Join(query, " AND ") + " "
	sql += "ORDER BY SO.IDOrder ASC "
	sql += "FOR UPDATE "

	if rs, err := tx.Query(sql, param...); err != nil {
		*errArrPtr = append(*errArrPtr, err)
	} else {
		defer rs.Close()

		var sqlField = struct {
			IDOrder int
		}{}
		sqlFieldArrPtr := util.StrutToSliceOfFieldAddress(&sqlField)

		for rs.Next() {
			if err := rs.Scan(sqlFieldArrPtr...); err != nil {
				*errArrPtr = append(*errArrPtr, err)
				break
			}

			outArr = append(outArr, sqlField.IDOrder)
		}

		if err := rs.Err(); err != nil {
			*errArrPtr = append(*errArrPtr, err)
		}
	}

	return outArr
}
