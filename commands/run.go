package commands

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a binary.",
	Long:  `run starts a binary if there's an actual payload.`,
	Run:   startRun,
}

func startRun(cmd *cobra.Command, args []string) {
	checkFlags()
	if Exec != "" {
		stdin := readStdin()
		if stdin != "" {
			log.Print(fmt.Sprintf("stdin='%s' exec='%s'", strings.TrimSpace(stdin), Exec))
			runCommand(Exec)
		} else {
			log.Print(fmt.Sprintf("stdin was blank"))
		}
	}

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

func checkFlags() {
	if Exec == "" {
		fmt.Println("Need a command to exec with '-e'")
		os.Exit(0)
	}
}

func readStdin() string {
	bytes, _ := ioutil.ReadAll(os.Stdin)
	stdin := string(bytes)
	if stdin == "" || stdin == "\n" || stdin == "[]" {
		return ""
	} else {
		return stdin
	}
}

func init() {
	RootCmd.AddCommand(runCmd)
}
