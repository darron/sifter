package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"time"
)

var eventCmd = &cobra.Command{
	Use:     "event",
	Short:   "Run a binary from a Consul event watch.",
	Aliases: []string{"run"},
	PreRun: func(cmd *cobra.Command, args []string) {
		checkEventFlags()
	},
	Long: `event runs a binary if there's an actually new Consul event.`,
	Run:  startEvent,
}

func startEvent(cmd *cobra.Command, args []string) {
	var oldEvent int64
	start := time.Now()

	stdin := readStdin()
	if stdin != "" {
		EventName, lTime, Payload := decodeStdin(stdin)
		lTimeString := strconv.FormatInt(int64(lTime), 10)
		ConsulKey := createKey(EventName)

		c, _ := Connect()
		ConsulData := Get(c, ConsulKey)
		if ConsulData != "" {
			oldEvent, _ = strconv.ParseInt(ConsulData, 10, 64)
		}

		if ConsulData == "" || oldEvent < lTime {
			Set(c, ConsulKey, lTimeString)
			runCommand(Exec, Payload)
			RunTime(start, "complete", fmt.Sprintf("watch='event' exec='%s' ltime='%d'", Exec, lTime))
			if DogStatsd {
				StatsdRunTime(start, EventName, Exec, lTime)
			}
		} else {
			RunTime(start, "duplicate", fmt.Sprintf("watch='event' exec='%s' ltime='%d'", Exec, lTime))
		}

	} else {
		RunTime(start, "blank", fmt.Sprintf("watch='event' exec='%s'", Exec))
	}

}

func checkEventFlags() {
	if Exec == "" {
		fmt.Println("Need a command to exec with '-e'")
		os.Exit(0)
	}
}

func init() {
	RootCmd.AddCommand(eventCmd)
}
