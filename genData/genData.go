package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)
import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func init() {
	rand.Seed(time.Now().UnixNano())

	//Reader = &devReader{name: "/dev/urandom"}
}

func RandomNumber(min, max int) int {
	return rand.Intn(max-min) + min
}

type Bug struct {
	BugID     int    `db:"BugID"`
	ProjectID int    `db:"ProjectID"`
	Status    int    `db:"Status"`
	Summary   string `db:"Summary"`
	Created   string `db:"Created"`
}

type BugBody struct {
	BugID int    `db:"BugID"`
	Body  string `db:"Body"`
}

func main() {
	var db *sqlx.DB
	var err error

	db, err = sqlx.Connect("mysql", "exp:exp@tcp(127.0.0.1:3306)/exp")

	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	defer db.Close()

	//
	var ns1 *sqlx.NamedStmt

	ns1, err = db.PrepareNamed(`INSERT INTO Bug (ProjectID, Status, Summary, Created) VALUES (:ProjectID, :Status, :Summary, NOW())`)

	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	//
	var ns2 *sqlx.NamedStmt

	ns2, err = db.PrepareNamed(`INSERT INTO BugBody (BugID, Body) VALUES (:BugID, :Body)`)

	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	//
	var rs sql.Result
	var lastInsertID int64

	bug := Bug{}
	bugBody := BugBody{}
	wordArr := []string{
		"aaaaaaaaa",
		"bbbbbbbbb",
		"ccccccccc",
		"ddddddddd",
		"eeeeeeeee",
		"fffffffff",
		"ggggggggg",
		"hhhhhhhhh",
		"iiiiiiiii",
		"jjjjjjjjj",
	}
	wordArrLen := len(wordArr)
	str := strings.Repeat("test ", 1000)
	tagStr := ""

	for i := 0; i < 1000000; i++ {
		bug.ProjectID = RandomNumber(0, wordArrLen)
		bug.Status = RandomNumber(0, wordArrLen)

		tagStr = wordArr[bug.ProjectID] + " " + wordArr[bug.Status] + " "

		bug.Summary = tagStr + strconv.Itoa(i)

		if rs, err = ns1.Exec(bug); err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			return
		}

		if lastInsertID, err = rs.LastInsertId(); err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			return
		}

		bugBody.BugID = int(lastInsertID)
		bugBody.Body = tagStr + str + strconv.Itoa(int(lastInsertID))

		if _, err = ns2.Exec(bugBody); err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			return
		}

		fmt.Printf("%d ", i)
	}

}
