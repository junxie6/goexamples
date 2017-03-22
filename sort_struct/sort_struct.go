package main

import (
	"fmt"
	"sort"
	"strings"
)

// Person
type Person struct {
	Name string
	Age  int
}

// PersonSortByName
type PersonSortByName []Person

func (p PersonSortByName) Len() int {
	return len(p)
}

func (p PersonSortByName) Less(j int, k int) bool {
	// sort by name only
	return p[j].Name < p[k].Name
}

func (p PersonSortByName) Swap(j int, k int) {
	p[j], p[k] = p[k], p[j]
}

// PersonSortByNameAge
type PersonSortByNameAge []Person

func (p PersonSortByNameAge) Len() int {
	return len(p)
}

func (p PersonSortByNameAge) Less(j int, k int) bool {
	// sort by name, then age
	if p[j].Name < p[k].Name {
		return true
	}
	if p[j].Name > p[k].Name {
		return false
	}
	return p[j].Age < p[k].Age
}

func (p PersonSortByNameAge) Swap(j int, k int) {
	p[j], p[k] = p[k], p[j]
}

//  Child
type Child struct {
	Name string
	Age  int
}

func main() {
	// Old way - implment the sort interface.
	p := []Person{
		Person{Name: "aaa", Age: 22},
		Person{Name: "bbb", Age: 19},
		Person{Name: "aaa", Age: 23},
		Person{Name: "aaa", Age: 100},
		Person{Name: "ccc", Age: 19},
		Person{Name: "aaa", Age: 19},
		Person{Name: "aaa", Age: 21},
		Person{Name: "aaa", Age: 20},
		Person{Name: "bbb", Age: 20},
		Person{Name: "ccc", Age: 21},
		Person{Name: "bbb", Age: 21},
		Person{Name: "ccc", Age: 20},
	}

	fmt.Printf("Old sorting method (ASC) - sort by name:\n")
	sort.Sort(PersonSortByName(p))
	outputPeople(p)

	fmt.Printf("Old sorting method (DESC) - sort by name:\n")
	sort.Sort(sort.Reverse(PersonSortByName(p)))
	outputPeople(p)

	fmt.Printf("Old sorting method (ASC) - sort by name then age:\n")
	sort.Sort(PersonSortByNameAge(p))
	outputPeople(p)

	// New way as of Go 1.8
	children := []Child{
		Child{Name: "aaa", Age: 22},
		Child{Name: "bbb", Age: 19},
		Child{Name: "aaa", Age: 23},
		Child{Name: "aaa", Age: 100},
		Child{Name: "ccc", Age: 19},
		Child{Name: "aaa", Age: 19},
		Child{Name: "aaa", Age: 21},
		Child{Name: "aaa", Age: 20},
		Child{Name: "bbb", Age: 20},
		Child{Name: "ccc", Age: 21},
		Child{Name: "bbb", Age: 21},
		Child{Name: "ccc", Age: 20},
	}

	fmt.Printf("New sorting method (ASC) - sort by name then age:\n")
	sort.Slice(children, func(i, j int) bool {
		switch strings.Compare(children[i].Name, children[j].Name) {
		case -1:
			return true
		case 1:
			return false
		}
		return children[i].Age < children[j].Age
	})

	outputChildren(children)
}

func outputPeople(p []Person) {
	for _, v := range p {
		fmt.Printf("%s %d\n", v.Name, v.Age)
	}
}

func outputChildren(p []Child) {
	for _, v := range p {
		fmt.Printf("%s %d\n", v.Name, v.Age)
	}
}
