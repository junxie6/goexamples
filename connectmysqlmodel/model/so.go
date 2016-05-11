package model

import (
	"log"
)

// SO is SalesOrder
type SO struct {
	IDOrder int
	Status  int
}

func (so *SO) Insert() {
	so.Status = 1

	stmt, err := db.Prepare("INSERT INTO SO (Status, changed) VALUES (?, NOW())")

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
