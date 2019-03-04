// Reference:
// https://appliedgo.net/mapreduce/
// https://play.golang.org/p/cipGuzMNT3
package main

import (
	"fmt"
)

func length(s string) int {
	return len(s)
}

func apply(list []string, fn func(string) int) []int {
	res := make([]int, len(list))
	for i, elem := range list {
		res[i] = fn(elem)
	}
	return res
}

func reduce(list []int, fn func(int, int) int) (res int) {
	for _, elem := range list {
		res = fn(res, elem)
	}
	return res
}

func sum(a, b int) int {
	return a + b
}

func main() {
	list := []string{"a", "bcd", "ef", "g", "hij"}
	res := reduce(apply(list, length), sum)
	fmt.Println(res)
}
