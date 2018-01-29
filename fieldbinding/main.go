package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
)

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/junxie6/goexamples/fieldbinding/fieldbinding"
)

var (
	db *sql.DB
)

// Table definition
// CREATE TABLE `salorder` (
//   `IDOrder` int(10) unsigned NOT NULL AUTO_INCREMENT,
//   `IsClose` tinyint(4) NOT NULL,
//   `IsConfirm` tinyint(4) NOT NULL,
//   `IDUser` int(11) NOT NULL,
//   `Created` datetime NOT NULL,
//   `Changed` datetime NOT NULL,
//   PRIMARY KEY (`IDOrder`),
//   KEY `IsClose` (`IsClose`)
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

func main() {
	var err error

	// starting database server
	db, err = sql.Open("mysql", "Username:Password@tcp(Host:Port)/DBName?parseTime=true")

	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}

	defer db.Close()

	// SampleQuery
	if v, err := SampleQuery(); err != nil {
		fmt.Printf("%s\n", err.Error())
	} else {
		var b bytes.Buffer

		if err := json.NewEncoder(&b).Encode(v); err != nil {
			fmt.Printf("SampleQuery: %v\n", err.Error())
		}

		fmt.Printf("SampleQuery: %v\n", b.String())
	}
}

func SampleQuery() ([]interface{}, error) {
	param := []interface{}{}

	param = append(param, 1)

	sql := "SELECT "
	sql += "  SalOrder.IDOrder "
	sql += ", SalOrder.IsClose "
	sql += ", SalOrder.IsConfirm "
	sql += ", SalOrder.IDUser "
	sql += ", SalOrder.Created "
	sql += "FROM SalOrder "
	sql += "WHERE "
	sql += "IsConfirm = ? "
	sql += "ORDER BY SalOrder.IDOrder ASC "

	rs, err := db.Query(sql, param...)

	if err != nil {
		return nil, err
	}

	defer rs.Close()

	// create a fieldbinding object.
	var fArr []string
	fb := fieldbinding.NewFieldBinding()

	if fArr, err = rs.Columns(); err != nil {
		return nil, err
	}

	fb.PutFields(fArr)

	//
	outArr := []interface{}{}

	for rs.Next() {
		if err := rs.Scan(fb.GetFieldPtrArr()...); err != nil {
			return nil, err
		}

		fmt.Printf("Row: %v, %v, %v, %s\n", fb.Get("IDOrder"), fb.Get("IsConfirm"), fb.Get("IDUser"), fb.Get("Created"))
		outArr = append(outArr, fb.GetFieldArr())
	}

	if err := rs.Err(); err != nil {
		return nil, err
	}

	return outArr, nil
}
