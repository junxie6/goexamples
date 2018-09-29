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
	"github.com/junxie6/util"
)

func init() {
	rand.Seed(time.Now().UnixNano())

	//Reader = &devReader{name: "/dev/urandom"}
}

type Bug struct {
	BugID     int     `db:"BugID"`
	ProjectID int     `db:"ProjectID"`
	Status    int     `db:"Status"`
	Amount    float64 `db:"Amount"`
	Qty       int     `db:"Qty"`
	Summary   string  `db:"Summary"`
	Created   string  `db:"Created"`
	Year      int8    `db:"Year"`
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

	ns1, err = db.PrepareNamed(`INSERT INTO Bug (ProjectID, Status, Amount, Qty, Summary, Created, Year) VALUES (:ProjectID, :Status, :Amount, :Qty, :Summary, :Created, :Year)`)

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

	for i := 0; i < 1; i++ {
		bug.ProjectID = util.RandomNumber(0, wordArrLen)
		bug.Status = util.RandomNumber(0, wordArrLen)

		tagStr = wordArr[bug.ProjectID] + " " + wordArr[bug.Status] + " "

		bug.Amount = float64(bug.ProjectID)
		bug.Qty = bug.ProjectID
		bug.Summary = tagStr + strconv.Itoa(i)

		bug.Created = util.RandomDate(3)
		bug.Year = bug.Created[2:4]

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
