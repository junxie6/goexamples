package gomap

import (
	"fmt"
	"testing"
)

func TestMain(t *testing.T) {
	// t.Errorf("err: %v", err)

	var str string

	m := NewGoMap()

	m.Add("asdf", "Hello World 1")
	str = m.Get("asdf")
	fmt.Printf("%v\n", str)

	m.Add("asdf", "Hello World 2")
	str = m.Get("asdf")
	fmt.Printf("%v\n", str)
}
