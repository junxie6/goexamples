package main

// Reference:
// https://gobyexample.com/signals
// https://stackoverflow.com/questions/11268943/is-it-possible-to-capture-a-ctrlc-signal-and-run-a-cleanup-function-in-a-defe?rq=1
// https://stackoverflow.com/questions/13107958/what-exactly-does-runtime-gosched-do

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time" // or "runtime"
)

func cleanup(exitStatus chan int) {
	fmt.Println("Please wait. Doing some cleanup tasks 1/3")
	time.Sleep(1 * time.Second)
	fmt.Println("Please wait. Doing some cleanup tasks 2/3")
	time.Sleep(1 * time.Second)
	fmt.Println("Please wait. Doing some cleanup tasks 3/3")
	time.Sleep(1 * time.Second)

	// TODO: implement the logic to determine the exit status
	exitStatus <- 0
}

func main() {
	go func() {
		// Make a buffered channel of size one, so the notifier will not be blocked
		sigs := make(chan os.Signal, 1)
		exitStatus := make(chan int, 1)

		// Register the given channel to receive notifications of the specified signals.
		signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

		// awaiting the signal
		<-sigs

		// Catch and handle Unix signals
		go cleanup(exitStatus)

		// Give the process a chance to cleanup itself.
		select {
		case <-time.After(10 * time.Second):
			fmt.Printf("Time is up. Exitting now.\n")
			os.Exit(1)
		case code := <-exitStatus:
			fmt.Printf("Exit status: %v\n", code)
			os.Exit(code)
		}
	}()

	for {
		fmt.Println("sleeping...")
		time.Sleep(10 * time.Second) // or runtime.Gosched() or similar per @misterbee
	}
}
