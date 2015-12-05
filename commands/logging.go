package commands

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Log sends logs to syslog.
func Log(message, priority string) {
	switch {
	case priority == "debug":
		if os.Getenv("SIFTER_DEBUG") != "" {
			log.Print(message)
		}
	default:
		log.Print(message)
	}

}

// RunTime shows logs with timing information.
func RunTime(start time.Time, location, extra string) {
	elapsed := time.Since(start)
	log := fmt.Sprintf("location='%s' elapsed='%s'", location, elapsed)
	if extra != "" {
		log = log + fmt.Sprintf(" %s", extra)
	}
	Log(log, "info")
}
