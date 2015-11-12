package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

func runCommand(command string) bool {
	parts := strings.Fields(command)
	cli := parts[0]
	args := parts[1:len(parts)]
	cmd := exec.Command(cli, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println("exec='error' message='%v'", err)
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

func decodeStdin(data string) (string, int64) {
	var events ConsulEvent
	err := json.Unmarshal([]byte(data), &events)
	if err != nil {
		fmt.Println("%#v", err)
	}
	name := string(events.Name)
	lTime := int64(events.LTime)
	log.Print(fmt.Sprintf("decoded event='%s' ltime='%d'", name, lTime))
	return name, lTime
}
