package main

// https://groups.google.com/forum/#!topic/golang-nuts/-eo7navkp10
// https://play.golang.org/p/NIEzTx8sjP
import (
	"errors"
	"os"
)

func DoStuff() (s string, err error) {
	myFile, err := os.Open("filename")
	if err != nil {
		return "", errors.New("Couldn't open file: " + err.Error())
	}

	defer func() {
		// Preferred by me version
		if cerr := myFile.Close(); cerr != nil && err == nil {
			err = cerr
		}

		// Alternative which I don't use
		err = errOr(err, myFile.Close())
	}()

	// Do stuff
	return "someString", nil
}

// provides ||-like operation for error values
func errOr(e, f error) error {
	if e == nil {
		e = f
	}
	return e
}

func main() {}
