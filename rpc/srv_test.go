package srv_test

import (
	"log"
	"net/rpc"
	"testing"
)

func BenchmarkConnection1(b *testing.B) {
	conn, err := rpc.Dial("tcp", "localhost:8080")

	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		line := []byte("HI")
		reply := new([]byte)

		err := conn.Call("Myfunc.Echo", line, reply)

		// If an error is returned, the reply parameter will not be sent back to the client.
		if err != nil {
			log.Printf("Call: %#v\n", err)
		}
	}

	conn.Close()
}

func BenchmarkConnection2(b *testing.B) {
	jobs := make(chan bool, 0)

	// This starts up 3 workers, initially blocked
	// because there are no jobs yet.
	for w := 1; w <= 500; w++ {
		go worker(jobs)
	}

	for i := 0; i < b.N; i++ {
		jobs <- true
	}

	close(jobs)
}

func worker(jobs <-chan bool) {
	for range jobs {
		conn, err := rpc.Dial("tcp", "localhost:8080")

		if err != nil {
			log.Fatal(err)
		}

		line := []byte("HI")
		reply := new([]byte)

		err = conn.Call("Myfunc.Echo", line, reply)

		// If an error is returned, the reply parameter will not be sent back to the client.
		if err != nil {
			log.Printf("Call: %#v\n", err)
		}

		conn.Close()
	}
}
