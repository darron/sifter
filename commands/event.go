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
	var previousLTime int64
	start := time.Now()

	stdin := readStdin()
	if stdin != "" {
		watchEvent := decodeEventStdin(stdin)
		watchEvent.examine()
		// Get event lTime
		lTime := watchEvent.getLTime()
		lTimeString := watchEvent.getLTimeString()
		// Get the event Name.
		eventName := watchEvent.getEventName()
		// Grab the payload - if any.
		payload := watchEvent.getPayload()
		// Grab the URL we will use to check Consul's previous info.
		nodeUrl := watchEvent.makeURL()
		// Connect to Consul
		consul, _ := Connect()
		// Get the previous value from Consul.
		previousData := Get(consul, nodeUrl)
		// If there's previousData - turn it into an int64.
		if previousData != "" {
			previousLTime, _ = strconv.ParseInt(previousData, 10, 64)
		}
		// If there's no previousData OR the previousLTime is less than the current lTime.
		// Then it's a new event - let's do the thing.
		if previousData == "" || previousLTime < lTime {
			Set(consul, nodeUrl, lTimeString)
			runCommand(Exec, payload)
			RunTime(start, "complete", fmt.Sprintf("watch='event' exec='%s' ltime='%d'", Exec, lTime))
			StatsdRunTime(start, Exec, "event", eventName, lTimeString)
		} else {
			RunTime(start, "duplicate", fmt.Sprintf("watch='event' exec='%s' ltime='%d'", Exec, lTime))
			StatsdDuplicate("event", eventName)
		}
	} else {
		RunTime(start, "blank", fmt.Sprintf("watch='event' exec='%s'", Exec))
		StatsdBlank("event")
	}
}

func checkEventFlags() {
	if Exec == "" {
		fmt.Println("Need a command to exec with '-e'")
		os.Exit(0)
	}
}

type EventWatch struct {
	Id            string `json:"ID"`
	Name          string `json:"Name"`
	Payload       string `json:"Payload,omitempty"`
	NodeFilter    string `json:"NodeFilter,omitempty"`
	ServiceFilter string `json:"ServiceFilter"`
	TagFilter     string `json:"TagFilter"`
	Version       int    `json:"Version"`
	LTime         int    `json:"LTime"`
}

func decodeEventStdin(data string) *EventWatch {
	events := make([]EventWatch, 0)
	err := json.Unmarshal([]byte(data), &events)
	if err != nil {
		Log(fmt.Sprintf("error: %s", data), "info")
		os.Exit(1)
	}
	sortutil.DescByField(events, "LTime")
	event := events[0]
	return &event
}

func (w *EventWatch) makeURL() string {
	hostname := getHostname()
	eventName := w.getEventName()
	url := fmt.Sprintf("%s/event/%s/%s", Prefix, eventName, hostname)
	return url
}

func (w *EventWatch) getEventName() string {
	name := w.Name
	return name
}

func (w *EventWatch) getPayload() string {
	payload := w.Payload
	return payload
}

func (w *EventWatch) getLTime() int64 {
	lTime := int64(w.LTime)
	return lTime
}

func (w *EventWatch) getLTimeString() string {
	lTime := w.getLTime()
	lTimeString := strconv.FormatInt(lTime, 10)
	return lTimeString
}

func (w *EventWatch) examine() {
	name := w.getEventName()
	lTime := w.getLTime()
	payload := w.getPayload()
	Log(fmt.Sprintf("decoded event='%s' ltime='%d' payload='%s'", name, lTime, payload), "debug")
}

func init() {
	RootCmd.AddCommand(eventCmd)
}
