package commands

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func runCommand(command, payload string) bool {
	parts := strings.Fields(command)
	cli := parts[0]
	args := parts[1:len(parts)]
	if payload != "" {
		args = append(args, payload)
		Log(fmt.Sprintf("exec='payload' payload='%s'", payload), "info")
	}
	cmd := exec.Command(cli, args...)
	Log(fmt.Sprintf("exec='runCommand' cli='%s' args='%s'", cli, args), "debug")
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
