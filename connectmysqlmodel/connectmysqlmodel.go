package main

import (
	"github.com/junhsieh/goexamples/connectmysqlmodel/model"
	"log"
)

const (
	dataSourceName = "go_erp:go_erp@tcp(localhost:3306)/go_erp"
)

func main() {
	if db, err := model.InitDB(dataSourceName); err != nil {
		log.Printf("DB Error: %v", err)
	} else {
		defer db.Close()

		so := model.SO{}
		so.EditSO()
	}
}
