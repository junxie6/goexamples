package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
)

const (
	exp1TplHTML = `
	<p>Host: {{.Host}}</p>
	<p>Data:	{{.Data}}</p>
`
)

// This example shows:
// 1. A shortcut way to create a new HTML template.
// 2.	Read HTML from a string
func serveExample1(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var v = struct {
		Host string
		Data string
	}{
		r.Host,
		"Test",
	}

	exp1Tpl := template.Must(template.New("").Parse(exp1TplHTML))
	exp1Tpl.Execute(w, v)
}

// This example shows:
// 1. Read HTML from a file.
// 2. Writing the output to a string instead of io.Writer.
func serveExample2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var v = struct {
		MyMsg string
	}{
		"Hello World",
	}

	exp2TplPath := "htmltemplate/template/test.tpl.html"
	exp2Tpl := template.New(path.Base(exp2TplPath)) // the template name needs to be the filename.

	var outStr string

	if exp2Tpl, err := exp2Tpl.ParseFiles(exp2TplPath); err != nil {
		log.Println(err)
	} else {
		var doc bytes.Buffer

		exp2Tpl.Execute(&doc, v)
		outStr = doc.String()
	}

	w.Write([]byte(outStr))
}

func main() {
	http.HandleFunc("/ex1", serveExample1)
	http.HandleFunc("/ex2", serveExample2)

	err := http.ListenAndServe(":80", nil)

	if err != nil {
		fmt.Printf("main(): %s\n", err)
		log.Fatal("ListenAndServe: ", err)
	}
}
