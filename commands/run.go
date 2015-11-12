package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a binary.",
	Long:  `run starts a binary if there's an actually new Consul event.`,
	Run:   startRun,
}

func startRun(cmd *cobra.Command, args []string) {
	var oldEvent int64
	checkFlags()

	if Exec != "" {
		stdin := readStdin()
		if stdin != "" {
			EventName, lTime := decodeStdin(stdin)
			log.Print(fmt.Sprintf("decoded event='%s' ltime='%d'", EventName, lTime))
			ConsulKey := createKey(EventName)

			c, _ := Connect()
			ConsulData := Get(c, ConsulKey)
			if ConsulData != "" {
				oldEvent, _ = strconv.ParseInt(ConsulData, 10, 64)
			}

			if ConsulData == "" || oldEvent < lTime {
				intString := strconv.FormatInt(int64(lTime), 10)
				Set(c, ConsulKey, intString)
				log.Print(fmt.Sprintf("ltime='%s' exec='%s'", intString, Exec))
				runCommand(Exec)
			} else {
				log.Print(fmt.Sprintf("status='old' stopping"))
			}

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
