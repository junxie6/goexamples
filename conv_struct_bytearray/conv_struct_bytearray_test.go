package conv_struct_bytearray_test

import (
	"encoding/gob"
	//"fmt"
	"github.com/oxtoacart/bpool"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

var (
	gBuf *bpool.BufferPool
)

type Person struct {
	Name string
	Age  int
}

func init() {
	gBuf = bpool.NewBufferPool(10)
}

// go test -run=^$ -bench="EncodeGoB"
func BenchmarkEncodeGoB(b *testing.B) {
	N := int32(b.N)
	procs := runtime.NumCPU()

	//
	var wg sync.WaitGroup
	wg.Add(procs)
	b.StartTimer()

	for p := 0; p < procs; p++ {
		go func() {
			for atomic.AddInt32(&N, -1) >= 0 {
				workerGoB()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func workerGoB() {
	p1 := Person{Name: "Jun", Age: 19}
	p2 := Person{}

	buf := gBuf.Get()

	gob.NewEncoder(buf).Encode(&p1)
	gob.NewDecoder(buf).Decode(&p2)

	//fmt.Printf("HERE %#v\n", p2)

	gBuf.Put(buf)
}
