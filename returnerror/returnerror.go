package main

import (
	"errors"
	"fmt"
)

func example1() error {
	errAnother := errors.New("another error")

	return errors.New("test 1: " + errAnother.Error())
}

func example2() error {
	return fmt.Errorf("test 2")
}

func main() {
	// first way, using errors package
	if err := example1(); err != nil {
		fmt.Printf("err: %v\n", err)
	}

	// second way, using fmt package
	if err := example2(); err != nil {
		fmt.Printf("err: %v\n", err)
	}
}
