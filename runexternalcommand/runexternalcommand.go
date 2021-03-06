package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

func main() {
	dumpDir := ""
	args := []string{"-h", "localhost", "-u", "DBUser", "-p", "-P", "3306", "DBName"}

	// Open a file
	out, err := os.OpenFile(path.Join(dumpDir, "dump.sql"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)

	if err != nil {
		fmt.Printf("%v", err)
	}

	defer out.Close()

	// Run command
	cmd := exec.Command("/usr/bin/mysqldump", args...)
	cmd.Stdout = out

	// Alternative Stdout
	//var out bytes.Buffer
	//cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		fmt.Printf("%v", err)
	}

	//if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
	// fmt.Printf("%v", err)
	//}
}
