package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

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

func runCommand(command, payload string) bool {
	parts := strings.Fields(command)
	cli := parts[0]
	args := parts[1:len(parts)]
	if payload != "" {
		args = append(args, payload)
		Log(fmt.Sprintf("exec='payload' payload='%s'", payload), "info")
	}
	cmd := exec.Command(cli, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		Log(fmt.Sprintf("exec='error' message='%v'", err), "info")
		return false
	} else {
		return true
	}
}

func getHostname() string {
	hostname, _ := os.Hostname()
	return hostname
}

func createKey(event string) string {
	hostname := getHostname()
	return fmt.Sprintf("sifter/%s/%s", event, hostname)
}

func readStdin() string {
	bytes, _ := ioutil.ReadAll(os.Stdin)
	stdin := string(bytes)
	if stdin == "" || stdin == "[]\n" || stdin == "\n" {
		return ""
	} else {
		return stdin
	}
}

func decodeStdin(data string) (string, int64, string) {
	events := make([]ConsulEvent, 0)
	err := json.Unmarshal([]byte(data), &events)
	if err != nil {
		Log(fmt.Sprintf("error: %s", data), "info")
		os.Exit(1)
	}
	var name = ""
	var lTime = int64(0)
	var payload = ""
	for _, event := range events {
		name = string(event.Name)
		if int64(event.LTime) > lTime {
			lTime = int64(event.LTime)
		}
		payload = event.Payload
	}
	Log(fmt.Sprintf("decoded event='%s' ltime='%d'", name, lTime), "info")
	return name, lTime, payload
}
