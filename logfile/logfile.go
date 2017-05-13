package main

import (
	"fmt"
	"log"
	"os"
	"syscall"
)

func logfile_example1(fileName string) (*os.File, error) {
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)

	if err != nil {
		return nil, err
	}

	// Set output to file.
	log.SetOutput(f)

	//log.SetOutput(os.Stderr) // to standard error (default).
	//log.SetOutput(os.Stdout) // to standard output.

	// Optional: include line number where the log is called.
	log.SetFlags(log.LstdFlags | log.Llongfile) // or log.Lshortfile.

	return f, nil
}

func logfile_example2(fileName string) (*os.File, error) {
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)

	if err != nil {
		return nil, err
	}

	syscall.Dup2(int(f.Fd()), 1) /* -- stdout */
	syscall.Dup2(int(f.Fd()), 2) /* -- stderr */

	// Optional: include line number where the log is called.
	log.SetFlags(log.LstdFlags | log.Llongfile) // or log.Lshortfile.

	return f, nil
}

// Reference:
// http://stackoverflow.com/questions/19965795/go-golang-write-log-to-file
// http://stackoverflow.com/questions/24809287/how-do-you-get-a-golang-program-to-print-the-line-number-of-the-error-it-just-ca
func main() {
	fileName := "test.log"

	// example 1
	//logFilePtr, err := logfile_example1(fileName)

	// example 2
	logFilePtr, err := logfile_example2(fileName)

	if err != nil {
		fmt.Printf("error opening file: %v", err)
		return
	}

	defer logFilePtr.Close()

	// log
	log.Println("This is a test log entry")
}
