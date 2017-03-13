// [Timers](timers) are for when you want to do
// something once in the future - _tickers_ are for when
// you want to do something repeatedly at regular
// intervals. Here's an example of a ticker that ticks
// periodically until we stop it.
//
// Note: Timers are for when you want to do something once in the future;
// tickers are for when you want to do something repeatedly at regular intervals.
//
// Reference: https://mmcgrana.github.io/2012/09/go-by-example-timers-and-tickers.html
package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Printf("Demoing ticker:\n")
	ticker()

	fmt.Printf("Demoing timer and ticker:\n")
	timer_and_ticker()
}

func ticker() {
	// Tickers use a similar mechanism to timers: a
	// channel that is sent values. Here we'll use the
	// `range` builtin on the channel to iterate over
	// the values as they arrive every 500ms.
	ticker := time.NewTicker(time.Millisecond * 500)

	go func() {
		for t := range ticker.C {
			fmt.Println("Tick at", t)
		}
	}()

	// Tickers can be stopped like timers. Once a ticker
	// is stopped it won't receive any more values on its
	// channel. We'll stop ours after 1600ms.
	time.Sleep(time.Millisecond * 1600)
	ticker.Stop()
	fmt.Println("Ticker stopped")
}

// A great feature of Go’s timers and tickers is that they hook into Go’s built-in concurrency mechanism: channels.
// This allows timers and tickers to interact seamlessly with other concurrent goroutines.
func timer_and_ticker() {
	timeChan := time.NewTimer(time.Second).C

	tickChan := time.NewTicker(time.Millisecond * 400).C

	doneChan := make(chan bool)
	go func() {
		time.Sleep(time.Second * 2)
		doneChan <- true
	}()

	for {
		select {
		case <-timeChan:
			fmt.Println("Timer expired")
		case <-tickChan:
			fmt.Println("Ticker ticked")
		case <-doneChan:
			fmt.Println("Done")
			return
		}
	}
}
