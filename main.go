package main

import (
	"./commands/"
	"log"
	"log/syslog"
	"runtime"
)

func main() {
	logwriter, e := syslog.New(syslog.LOG_NOTICE, "sifter")
	if e == nil {
		log.SetOutput(logwriter)
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
	commands.RootCmd.Execute()
}
