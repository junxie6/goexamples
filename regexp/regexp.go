package main

import (
	"fmt"
	"regexp"
)

func main() {
	var r *regexp.Regexp
	var err error

	str := "12345ms"
	expr := "^([0-9]+)(m?s)$"

	// method one
	if r, err = regexp.Compile(expr); err != nil {
		fmt.Printf("Err: %#v\n", err.Error())
	}

	resultArr := r.FindStringSubmatch(str)

	fmt.Printf("%#v\n", resultArr)

	// method two
	if ok, err := regexp.MatchString(expr, str); err != nil {
		fmt.Printf("Err: %#v\n", err.Error())
	} else if ok {
		fmt.Printf("It is matched.\n")
	}
}
