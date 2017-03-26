package conv_struct_bytearray_test

import (
	"bufio"
	"bytes"
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
func BenchmarkEncodeGoB1(b *testing.B) {
	N := int32(b.N)
	procs := runtime.NumCPU()

	//
	var wg sync.WaitGroup
	wg.Add(procs)
	b.StartTimer()

	for p := 0; p < procs; p++ {
		go func() {
			for atomic.AddInt32(&N, -1) >= 0 {
				encodeDecodeStruct()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkEncodeGoB2(b *testing.B) {
	var wg sync.WaitGroup
	jobs := make(chan bool, 4096)
	procs := runtime.NumCPU()

	// This starts up 3 workers, initially blocked
	// because there are no jobs yet.
	for w := 0; w < procs; w++ {
		wg.Add(1)

		go func() {
			worker(jobs)
			wg.Done()
		}()
	}

	for i := 0; i < b.N; i++ {
		jobs <- true
	}

	close(jobs)

	wg.Wait()
}

func worker(jobs <-chan bool) {
	for range jobs {
		encodeDecodeStruct()
	}
}

func encodeDecodeStruct() {
	p1 := Person{Name: "Jun", Age: 19}
	p2 := Person{}

	buf := gBuf.Get()

	gob.NewEncoder(buf).Encode(&p1)
	gob.NewDecoder(buf).Decode(&p2)

	//fmt.Printf("HERE %#v\n", p2)

	gBuf.Put(buf)
}

func BenchmarkEncodeGoBEncode1(b *testing.B) {
	buf := new(bytes.Buffer)
	p1 := Person{Name: "Jun", Age: 19}

	for i := 0; i < b.N; i++ {
		gob.NewEncoder(buf).Encode(p1)
	}
}

func BenchmarkEncodeGoBEncode2(b *testing.B) {
	buf := new(bytes.Buffer)
	en := gob.NewEncoder(buf)
	p1 := Person{Name: "Jun", Age: 19}

	for i := 0; i < b.N; i++ {
		en.Encode(p1)
	}
}

func BenchmarkEncodeGoBEncode3(b *testing.B) {
	bb := new(bytes.Buffer)
	buf := bufio.NewWriter(bb)
	en := gob.NewEncoder(buf)
	p1 := Person{Name: "Jun", Age: 19}

	for i := 0; i < b.N; i++ {
		en.Encode(p1)
		buf.Flush()
	}
}

func BenchmarkEncodeGoBDecode1(b *testing.B) {
	buf := new(bytes.Buffer)
	en := gob.NewEncoder(buf)
	p1 := Person{Name: "Jun", Age: 19}
	en.Encode(p1)

	for i := 0; i < b.N; i++ {
		gob.NewDecoder(buf).Decode(&p1)
	}
}

func BenchmarkEncodeGoBDecode2(b *testing.B) {
	buf := new(bytes.Buffer)
	en := gob.NewEncoder(buf)
	p1 := Person{Name: "Jun", Age: 19}
	en.Encode(p1)
	de := gob.NewDecoder(buf)

	for i := 0; i < b.N; i++ {
		de.Decode(&p1)
	}
}

func BenchmarkEncodeGoBDecode3(b *testing.B) {
	buf := new(bytes.Buffer)
	en := gob.NewEncoder(buf)
	p1 := Person{Name: "Jun", Age: 19}
	en.Encode(p1)
	de := gob.NewDecoder(bytes.NewReader(buf.Bytes()))

	for i := 0; i < b.N; i++ {
		de.Decode(&p1)
	}
}
