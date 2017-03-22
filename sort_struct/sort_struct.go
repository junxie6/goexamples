package main

import (
	"fmt"
	"sort"
)

type Person struct {
	Name string
	Age  int
}

type People []Person

func (p People) Len() int {
	return len(p)
}

func (p People) Less(j int, k int) bool {
	return p[j].Name < p[k].Name
}

func (p People) Swap(j int, k int) {
	p[j], p[k] = p[k], p[j]
}

func main() {
	p := People{
		Person{Name: "bbb", Age: 19},
		Person{Name: "ccc", Age: 19},
		Person{Name: "aaa", Age: 19},
	}

	sort.Sort(p)
	output(p)

	sort.Sort(sort.Reverse(p))
	output(p)
}

func output(p People) {
	for _, v := range p {
		fmt.Printf("%s %d\n", v.Name, v.Age)
	}
}
