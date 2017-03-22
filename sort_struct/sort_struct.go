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

type People []Person

func (p People) Len() int {
	return len(p)
}

func (p People) Less(j int, k int) bool {
	// sort by name only
	//return p[j].Name < p[k].Name

	// sort by name, then age
	if p[j].Name < p[k].Name {
		return true
	}
	if p[j].Name > p[k].Name {
		return true
	}
	return p[j].Age < p[k].Age
}

func (p People) Swap(j int, k int) {
	p[j], p[k] = p[k], p[j]
}

//  Child
type Child struct {
	Name string
	Age  int
}

func main() {
	// Old way - implment the sort interface.
	p := People{
		Person{Name: "aaa", Age: 22},
		Person{Name: "bbb", Age: 19},
		Person{Name: "aaa", Age: 23},
		Person{Name: "ccc", Age: 19},
		Person{Name: "aaa", Age: 19},
		Person{Name: "aaa", Age: 21},
		Person{Name: "aaa", Age: 20},
	}

	fmt.Printf("Old sorting method (ASC):\n")
	sort.Sort(p)
	outputPeople(p)

	fmt.Printf("Old sorting method (DESC):\n")
	sort.Sort(sort.Reverse(p))
	outputPeople(p)

	// New way as of Go 1.8
	children := []Child{
		Child{Name: "aaa", Age: 22},
		Child{Name: "bbb", Age: 19},
		Child{Name: "aaa", Age: 23},
		Child{Name: "ccc", Age: 19},
		Child{Name: "aaa", Age: 19},
		Child{Name: "aaa", Age: 21},
		Child{Name: "aaa", Age: 20},
	}

	fmt.Printf("New sorting method (ASC):\n")
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

func outputPeople(p People) {
	for _, v := range p {
		fmt.Printf("%s %d\n", v.Name, v.Age)
	}
}

func outputChildren(p []Child) {
	for _, v := range p {
		fmt.Printf("%s %d\n", v.Name, v.Age)
	}
}
