package main

import (
	"fmt"
	"github.com/darron/sifter/commands"
	"log"
	"log/syslog"
	"os"
	"runtime"
)

var (
	// CompileDate tracks when the binary was compiled. It's inserted during a build
	// with build flags. Take a look at the Makefile for information.
	CompileDate = "No date provided."

	// GitCommit tracks the SHA of the built binary. It's inserted during a build
	// with build flags. Take a look at the Makefile for information.
	GitCommit = "No revision provided."

	// Version is the version of the built binary. It's inserted during a build
	// with build flags. Take a look at the Makefile for information.
	Version = "No version provided."

	// GoVersion details the version of Go this was compiled with.
	GoVersion = runtime.Version()
)

func main() {
	logwriter, e := syslog.New(syslog.LOG_NOTICE, "sifter")
	if e == nil {
		log.SetOutput(logwriter)
	}

	args := os.Args[1:]
	for _, arg := range args {
		if arg == "-v" || arg == "--version" {
			fmt.Printf("%-8s : %s\n%-8s : %s\n%-8s : %s\n%-8s : %s\n",
				"Version", Version,
				"Revision", GitCommit,
				"Date", CompileDate,
				"Go", GoVersion)
			os.Exit(0)
		}
	}

	commands.RootCmd.Execute()
}
