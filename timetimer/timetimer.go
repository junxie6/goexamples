package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Printf("Demoing timer1:\n")
	timer1()

	fmt.Printf("Demoing timer2:\n")
	timer2()
}

// Timers represent a single event in the future. You tell the timer how long you want to wait,
// and it gives you a channel that will be notified at that time.
// Reference: https://mmcgrana.github.io/2012/09/go-by-example-timers-and-tickers.html
func timer1() {
	timer := time.NewTimer(time.Second * 2)
	<-timer.C

	println("Timer expired")
}

// If you just wanted to wait, you could have used time.Sleep.
// One reason a timer may be useful is that you can cancel the timer before it expires.
func timer2() {
	timer := time.NewTimer(time.Second)

	go func() {
		<-timer.C
		println("Timer expired")
	}()

	stop := timer.Stop()

	println("Timer cancelled:", stop)
}
