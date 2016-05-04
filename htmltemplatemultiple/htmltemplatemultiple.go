package main

import (
	"bytes"
	"fmt"
	"html/template"
)

func main() {
	var docHTML1 bytes.Buffer
	var docHTML2 bytes.Buffer

	var v = struct {
		MyMsgA string        // untrusted plain text
		MyMsgB template.HTML // trusted HTML
	}{
		MyMsgA: "<b>Hello World A</b>",
		MyMsgB: "<b>Hello World B</b>",
	}

	layoutTpl := "htmltemplatemultiple/template/testlayout.tpl.html"
	headerTpl := "htmltemplatemultiple/template/testheader.tpl.html"
	bodyTpl := "htmltemplatemultiple/template/testbody.tpl.html"
	footerTpl := "htmltemplatemultiple/template/testfooter.tpl.html"

	tmpl := template.Must(template.ParseFiles(layoutTpl, headerTpl, bodyTpl, footerTpl))

	fmt.Printf("Show the body region:\n")
	tmpl.ExecuteTemplate(&docHTML1, "body", v)
	str1 := docHTML1.String()
	fmt.Printf("%s\n", str1)

	fmt.Printf("Show the entire layout:\n")
	//tmpl.Execute(&docHTML2, v)
	tmpl.ExecuteTemplate(&docHTML2, "layout", v)
	str2 := docHTML2.String()
	fmt.Printf("%s\n", str2)
}
