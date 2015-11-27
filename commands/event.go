package commands

import (
	"encoding/json"
	"fmt"
	"github.com/pmylund/sortutil"
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
	Long: `event runs a binary if there's a new Consul event.`,
	Run:  startEvent,
}

func startEvent(cmd *cobra.Command, args []string) {
	var oldEvent int64
	start := time.Now()

	stdin := readStdin()
	if stdin != "" {
		EventName, lTime, Payload := decodeEventStdin(stdin)
		lTimeString := strconv.FormatInt(int64(lTime), 10)
		ConsulKey := createEventKey(EventName)

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

type ConsulEvent struct {
	Id            string `json:"ID"`
	Name          string `json:"Name"`
	Payload       string `json:"Payload,omitempty"`
	NodeFilter    string `json:"NodeFilter,omitempty"`
	ServiceFilter string `json:"ServiceFilter"`
	TagFilter     string `json:"TagFilter"`
	Version       int    `json:"Version"`
	LTime         int    `json:"LTime"`
}

func createEventKey(event string) string {
	hostname := getHostname()
	return fmt.Sprintf("sifter/event/%s/%s", event, hostname)
}

func decodeEventStdin(data string) (string, int64, string) {
	events := make([]ConsulEvent, 0)
	err := json.Unmarshal([]byte(data), &events)
	if err != nil {
		Log(fmt.Sprintf("error: %s", data), "info")
		os.Exit(1)
	}
	sortutil.DescByField(events, "LTime")
	event := events[0]
	name := event.Name
	lTime := int64(event.LTime)
	payload := event.Payload
	Log(fmt.Sprintf("decoded event='%s' ltime='%d' payload='%s'", name, lTime, payload), "info")
	return name, lTime, payload
}

func init() {
	RootCmd.AddCommand(eventCmd)
}
