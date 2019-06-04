package main

import (
	"fmt"
	"log"
	"log/syslog"
)

func main() {
	logWriter, err := syslog.Dial("udp", "192.168.6.66:9000", syslog.LOG_WARNING|syslog.LOG_DAEMON, "mytag")

	if err != nil {
		log.Fatal(err)
	}

	defer logWriter.Close()

	fmt.Fprintf(logWriter, "This is a daemon warning with mytag.")

	logWriter.Emerg("This is a daemon emergency with mytag.")
	logWriter.Alert("This is a daemon alert with mytag.")
	logWriter.Crit("This is a daemon critical with mytag.")
	logWriter.Err("This is a daemon error with mytag.")
	logWriter.Warning("This is a daemon warning with mytag.")
	logWriter.Notice("This is a daemon notice with mytag.")
	logWriter.Info("This is a daemon info with mytag.")
	logWriter.Debug("This is a daemon debug with mytag.")

	logWriter.Write([]byte("This is a daemon write with mytag."))
}
