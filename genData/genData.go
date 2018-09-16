package main

import (
	"database/sql"
	"fmt"
	"log"
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

	db, err = sqlx.Connect("mysql", "exp:exp@tcp(127.0.0.1:4000)/exp")

	if err != nil {
		log.Fatalln(err)
		return
	}

	defer db.Close()

	//
	var ns1 *sqlx.NamedStmt

	ns1, err = db.PrepareNamed(`INSERT INTO Bug (ProjectID, Status, Summary, Created) VALUES (:ProjectID, :Status, :Summary, NOW())`)

	if err != nil {
		log.Fatalln(err)
		return
	}

	//
	var ns2 *sqlx.NamedStmt

	ns2, err = db.PrepareNamed(`INSERT INTO BugBody (BugID, Body) VALUES (:BugID, :Body)`)

	if err != nil {
		log.Fatalln(err)
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

	for i := 0; i < 1000000; i++ {
		bug.ProjectID = RandomNumber(1, wordArrLen)
		bug.Status = RandomNumber(1, wordArrLen)
		bug.Summary = "test " + strconv.Itoa(i)

		if rs, err = ns1.Exec(bug); err != nil {
			log.Fatalln(err)
			return
		}

		if lastInsertID, err = rs.LastInsertId(); err != nil {
			log.Fatalln(err)
			return
		}

		bugBody.BugID = int(lastInsertID)
		bugBody.Body = wordArr[bug.ProjectID] + " " + wordArr[bug.Status] + " " + str + strconv.Itoa(int(lastInsertID))

		if _, err = ns2.Exec(bugBody); err != nil {
			log.Fatalln(err)
			return
		}

		fmt.Printf("%d ", i)
	}

}
