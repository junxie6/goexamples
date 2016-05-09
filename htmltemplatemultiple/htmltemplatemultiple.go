package main

import (
	"bytes"
	"fmt"
	"html/template"
)

func main() {
	var docHTML1 bytes.Buffer
	var docHTML2 bytes.Buffer

	var tplParam = struct {
		MyMsgA string        // untrusted plain text
		MyMsgB template.HTML // trusted HTML
	}{
		MyMsgA: "<b>Hello World A</b>",
		MyMsgB: "<b>Hello World B</b>",
	}

	layoutTpl := "htmltemplatemultiple/template/layout.tpl.html"
	headerTpl := "htmltemplatemultiple/template/header.tpl.html"
	bodyTpl := "htmltemplatemultiple/template/body.tpl.html"
	footerTpl := "htmltemplatemultiple/template/footer.tpl.html"

	tpl := template.Must(template.ParseFiles(layoutTpl, headerTpl, bodyTpl, footerTpl))

	fmt.Printf("Show the entire layout:\n")
	tpl.ExecuteTemplate(&docHTML1, "layout", tplParam)
	str1 := docHTML1.String()
	//bytes := buffer.Bytes() // Convert to []byte
	//bytes := buffer.WriteTo(w) // use Buffer.WriteTo() to copy the buffer contents directly to a Writer.
	fmt.Printf("%s\n", str1)

	fmt.Printf("Show the body region:\n")
	//tpl.Execute(&docHTML2, tplParam)
	tpl.ExecuteTemplate(&docHTML2, "body", tplParam)
	str2 := docHTML2.String()
	fmt.Printf("%s\n", str2)
}
