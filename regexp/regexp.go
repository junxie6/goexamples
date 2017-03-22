package main

import (
	"fmt"
	"regexp"
)

func main() {
	var r *regexp.Regexp
	var err error

	str := "12345ms"

	if r, err = regexp.Compile("^([0-9]+)(m?s)$"); err != nil {
		fmt.Printf("%#v\n", err.Error())
	}

	resultArr := r.FindStringSubmatch(str)

	fmt.Printf("%#v\n", resultArr)
}
