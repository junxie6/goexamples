package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"path"
)

func main() {
	var v = struct {
		MyMsg string
	}{
		"Hello World",
	}

	homeTplPath := "htmltemplate/template/test.tpl.html"
	homeTpl := template.New(path.Base(homeTplPath)) // the template name needs to be the filename.

	if homeTpl, err := homeTpl.ParseFiles(homeTplPath); err != nil {
		log.Println(err)
	} else {
		var doc bytes.Buffer

		homeTpl.Execute(&doc, v)
		str := doc.String()

		fmt.Printf("%s\n", str)
	}
}
