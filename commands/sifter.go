package commands

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

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

func readStdin() string {
	bytes, _ := ioutil.ReadAll(os.Stdin)
	stdin := string(bytes)
	if stdin == "" || stdin == "[]\n" || stdin == "\n" {
		return ""
	} else {
		return stdin
	}
}
