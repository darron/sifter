package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a binary.",
	Long:  `run starts a binary if there's an actually new Consul event.`,
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
			log.Print(fmt.Sprintf("stdin='blank' NOT running '%s'", Exec))
		}
	}

}

func checkFlags() {
	if Exec == "" {
		fmt.Println("Need a command to exec with '-e'")
		os.Exit(0)
	}
}

func init() {
	RootCmd.AddCommand(runCmd)
}
