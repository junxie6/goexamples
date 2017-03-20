// Reference:
// https://elithrar.github.io/article/using-buffer-pools-with-go/
// https://golang.org/pkg/sync/#Pool
package main

import (
	//"bytes"
	"encoding/gob"
	"fmt"
	"github.com/oxtoacart/bpool"
)

var (
	gBuf *bpool.BufferPool
)

type Person struct {
	Name string
	Age  int
}

type Car struct {
	Name  string
	Width int
}

func main() {
	gBuf = bpool.NewBufferPool(5)

	// In order to use our pool of workers we need to send
	// them work and collect their results. We make 2 channels for this.
	jobs := make(chan int, 100)
	results := make(chan bool, 100)
	numOfJobs := 10

	// This starts up 3 workers, initially blocked
	// because there are no jobs yet.
	for w := 0; w < 3; w++ {
		go worker(jobs, results)
	}

	for i := 0; i < numOfJobs; i++ {
		jobs <- i
	}

	close(jobs)

	// Finally we collect all the results of the work.
	for a := 0; a < numOfJobs; a++ {
		<-results
	}
}

func worker(jobs <-chan int, results chan<- bool) {
	for range jobs {
		p1 := Person{Name: "Jun", Age: 19}
		p2 := Person{}

		buf := gBuf.Get()

		gob.NewEncoder(buf).Encode(&p1)
		gob.NewDecoder(buf).Decode(&p2)

		fmt.Printf("HERE %#v\n", p2)

		gBuf.Put(buf)
		results <- true
	}
}
