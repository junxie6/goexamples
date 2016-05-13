package model

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
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
func (so *SO) editSO() []error {
	errArr := []error{}
	tx, err := db.Begin()

	if err != nil {
		return append(errArr, err)
	}

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

	return errArr
}

func (so *SO) lockSO(tx *sql.Tx, errArrPtr *[]error) {
}
