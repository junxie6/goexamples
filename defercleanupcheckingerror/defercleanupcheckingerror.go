package main

import (
	"io"
	"log"
)

type Resource struct {
	name string
}

func Open(name string) (*Resource, error) {
	return &Resource{name}, nil
}

func (r *Resource) Close() error {
	log.Printf("closing %s\n", r.name)
	return nil
}

func Close(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Printf(err.Error())
	}
}

// Reference:
// Deferred Cleanup, Checking Errors, and Potential Problems
// http://www.blevesearch.com/news/Deferred-Cleanup,-Checking-Errors,-and-Potential-Problems/
// Defer, Panic, and Recover
// http://blog.golang.org/defer-panic-and-recover
//
// "deferred function's arguments are evaluated when the defer statement is evaluated."
// "…they may refer to variables defined in a surrounding function. Those variables are then shared between the surrounding function and the function literal, and they survive as long as they are accessible."
func main() {
	// example one
	if r, err := Open("a"); err != nil {
		log.Printf(err.Error())
	} else {
		defer func(r *Resource) {
			if err := r.Close(); err != nil {
				log.Printf(err.Error())
			}
		}(r)
	}

	// example two (better)
	if r, err := Open("b"); err != nil {
		log.Printf(err.Error())
	} else {
		defer Close(r)
	}
}
