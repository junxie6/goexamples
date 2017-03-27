// Reference:
// https://groups.google.com/forum/#!topic/golang-nuts/n_By5xPzDho/discussion
// https://play.golang.org/p/RXFN2ACuNC
package main

import (
	"log"
	"math/rand"
	"os"
	//"runtime"
	"sync"
	"time"
)

type myPool struct {
	sync.Pool
}

func NewMyPool() *myPool {
	var m myPool
	m.Pool.New = func() interface{} {
		log.Println("Made a new buffer")
		return make([]byte, 10240)
	}
	return &m
}

func (m *myPool) Get() []byte {
	return m.Pool.Get().([]byte)
}

func (m *myPool) Put(b []byte) {
	// Reset buffer - first method
	// will reuse the already allocated capacity of the underlying array
	//b = b[:0]

	// Reset buffer - second method
	// will have to do reallocation while it grows.
	//b = []byte{}

	m.Pool.Put(b)
}

func main() {
	var p = NewMyPool()

	go func() {
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()

	f := func() {
		for {
			b := p.Get()
			time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
			p.Put(b)
			//runtime.GC()
		}
	}
	go f()
	go f()
	f()
}
