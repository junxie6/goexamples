package main

import (
	"fmt"
)

func test() string {
	return "hello world"
}

func InArrayInt(v int, vArr []int) bool {
	for _, vv := range vArr {
		if v == vv {
			return true
		}
	}

	return false
}

func main() {
	fmt.Printf(test())
}

// 任何 language，都逃不出以下幾個東西：
// 1: loop (迴圈)
// 2. if else (判斷)
// 3. variable (變數)
// 4. constant
// 5. function
