package main

import (
	"fmt"
)

type Mutatable struct {
	a int
	b int
}

func (m Mutatable) StayTheSame() {
	m.a = 3
	m.b = 3
}

func (m *Mutatable) Mutate() {
	m.a = 5
	m.b = 5
}

// A method call x.m() is valid if the method set of (the type of) x contains m and the argument list can be assigned to the parameter list of m. If x is addressable and &x’s method set contains m, x.m() is shorthand for (&x).m()
// Reference: http://nathanleclaire.com/blog/2014/08/09/dont-get-bitten-by-pointer-vs-non-pointer-method-receivers-in-golang/
func main() {
	m := Mutatable{9, 9}
	fmt.Println(m)

	m.StayTheSame()
	fmt.Println(m)

	(&m).StayTheSame()
	fmt.Println(m)

	m.Mutate()
	fmt.Println(m)

	(&m).Mutate()
	fmt.Println(m)
}
