package main

// An example streaming XML parser.
// Reference:
// https://github.com/dps/go-xml-parse/blob/master/go-xml-parse.go
// https://en.wikipedia.org/wiki/Wikipedia:Database_download
// https://dumps.wikimedia.org/enwiki/20190220/

import (
	"bufio"
	//"database/sql"
	"encoding/xml"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"net/url"
	"os"
	"regexp"
	//"strconv"
	"strings"
	"sync"
	//"time"
)

import (
	"github.com/jmoiron/sqlx"
)

var (
	DB        *sqlx.DB
	inputFile = flag.String("infile", "data/enwiki-20190220-pages-articles-multistream.xml", "Input file path")
	indexFile = flag.String("indexfile", "out/article_list.txt", "article list output file")
	filter, _ = regexp.Compile("^file:.*|^talk:.*|^special:.*|^wikipedia:.*|^wiktionary:.*|^user:.*|^user_talk:.*")
)

// Here is an example article from the Wikipedia XML dump
//
// <page>
// 	<title>Apollo 11</title>
//      <redirect title="Foo bar" />
// 	...
// 	<revision>
// 	...
// 	  <text xml:space="preserve">
// 	  {{Infobox Space mission
// 	  |mission_name=&lt;!--See above--&gt;
// 	  |insignia=Apollo_11_insignia.png
// 	...
// 	  </text>
// 	</revision>
// </page>
//
// Note how the tags on the fields of Page and Redirect below
// describe the XML schema structure.

type Redirect struct {
	Title string `xml:"title,attr"`
}

type Page struct {
	ID    string   `xml:"id"`
	Title string   `xml:"title"`
	Redir Redirect `xml:"redirect"`
	Text  string   `xml:"revision>text"`
}

func CanonicalizeTitle(title string) string {
	can := strings.ToLower(title)
	can = strings.Replace(can, " ", "_", -1)
	can = url.QueryEscape(can)
	return can
}

func FirstNChar(s string, n int) string {
	l := len(s)

	if l < n {
		return s[0:l]
	}
	return s[0:n]
}

func WritePage(title string, text string) {
	var outFile *os.File
	var err error

	if _, err := os.Stat("out/" + title); err == nil {
		fmt.Printf("Error: %s\n", "File exist: "+title)
		return
	}

	if outFile, err = os.Create("out/" + title); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	//
	writer := bufio.NewWriter(outFile)
	defer outFile.Close()

	writer.WriteString(text)
	writer.Flush()
}

func InitDB() error {
	var err error

	DB, err = sqlx.Connect("mysql", "exp:exp@tcp(127.0.0.1)/exp")

	if err != nil {
		return err
	}

	// TODO: move setting to config
	DB.SetMaxOpenConns(1000)
	DB.SetMaxIdleConns(300)
	DB.SetConnMaxLifetime(0)

	return err
}

func WritePageRecord(stmt *sqlx.NamedStmt, id string, title string, body string) error {
	var err error

	_, err = stmt.Exec(map[string]interface{}{
		"OrigID": id,
		"Title":  title,
		"Body":   body,
	})

	return err
}

func worker(id int, jobs <-chan Page, wg *sync.WaitGroup) {
	defer wg.Done()

	//
	var err error
	var stmt *sqlx.NamedStmt

	if _, err = DB.Exec("SET autocommit = 0"); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	if stmt, err = DB.PrepareNamed("INSERT INTO Article (OrigID, Title, Body, Data, Created) VALUES (:OrigID, :Title, :Body, '{}', NOW())"); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	defer stmt.Close()

	count := 0

	for j := range jobs {
		//fmt.Printf("Worker %d\n", id)
		//time.Sleep(5 * time.Second)

		//WritePage(j.ID+"_"+j.Title, j.Text)

		if err = WritePageRecord(stmt, j.ID, j.Title, j.Text); err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			continue
		}

		if id == 1 {
		}

		if (count % 10) == 0 {
			//fmt.Printf("%d Commit: %d\n", id, count)
			DB.Exec("COMMIT")
		}

		count++
	}

	// TODO:
	if (count % 1000) > 1 {
		//fmt.Printf("%d Committed: %d\n", id, count)
		DB.Exec("COMMIT")
	}
}

func main() {
	var err error
	var xmlFile *os.File

	flag.Parse()

	//
	if err = InitDB(); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	defer DB.Close()

	//

	//
	var wg sync.WaitGroup
	jobs := make(chan Page, 100)

	for w := 1; w <= 3; w++ {
		wg.Add(1)
		go worker(w, jobs, &wg)
	}

	//
	if xmlFile, err = os.Open(*inputFile); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	defer xmlFile.Close()

	//
	decoder := xml.NewDecoder(xmlFile)
	total := 0
	var inElement string
	var t xml.Token

	for {
		// DEBUG:
		if total == 99 {
			break
		}

		// Read tokens from the XML document in a stream.
		if t, err = decoder.Token(); err != nil {
			if err == io.EOF {
				break
			}

			fmt.Printf("Failed to read token: %s\n", err.Error())
			continue
			//panic("Failed to read token: " + err.Error())
		}

		if t == nil {
			break
		}

		// Inspect the type of the token just read.
		switch se := t.(type) {
		case xml.StartElement:
			// If we just read a StartElement token
			inElement = se.Name.Local

			// ...and its name is "page"
			if inElement == "page" {
				var p Page

				// decode a whole chunk of following XML into the
				// variable p which is a Page (se above)
				decoder.DecodeElement(&p, &se)

				// Do some stuff with the page.
				p.Title = FirstNChar(CanonicalizeTitle(p.Title), 16)
				m := filter.MatchString(p.Title)

				if !m && p.Redir.Title == "" {
					jobs <- p
					total++
				}
			}
		default:
		}

	}

	fmt.Printf("Closing jobs\n")
	close(jobs)

	fmt.Printf("Waiting jobs to finish\n")
	wg.Wait()

	fmt.Printf("Total articles: %d \n", total)
}
