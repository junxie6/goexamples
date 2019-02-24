package main

// An example streaming XML parser.
// Reference:
// https://github.com/dps/go-xml-parse/blob/master/go-xml-parse.go
// https://en.wikipedia.org/wiki/Wikipedia:Database_download
// https://dumps.wikimedia.org/enwiki/20190220/

import (
	"bufio"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"regexp"
	"strings"
)

var inputFile = flag.String("infile", "data/enwiki-20190220-pages-articles-multistream.xml", "Input file path")
var indexFile = flag.String("indexfile", "out/article_list.txt", "article list output file")

var filter, _ = regexp.Compile("^file:.*|^talk:.*|^special:.*|^wikipedia:.*|^wiktionary:.*|^user:.*|^user_talk:.*")

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

func WritePage(title string, text string) {
	var outFile *os.File
	var err error

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

func main() {
	var err error
	var xmlFile *os.File

	flag.Parse()

	//
	if xmlFile, err = os.Open(*inputFile); err != nil {
		fmt.Println("Error opening file:", err)
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
		if total == 10 {
			break
		}

		// Read tokens from the XML document in a stream.
		if t, err = decoder.Token(); err != nil {
			if err == io.EOF {
				break
			}

			panic("Failed to read token: " + err.Error())
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
				p.Title = CanonicalizeTitle(p.Title)
				m := filter.MatchString(p.Title)

				if !m && p.Redir.Title == "" {
					WritePage(p.Title, p.Text)
					total++
				}
			}
		default:
		}

	}

	fmt.Printf("Total articles: %d \n", total)
}
