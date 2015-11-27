package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var keyCmd = &cobra.Command{
	Use:     "key",
	Short:   "Run a binary from a Consul key watch.",
	PreRun: func(cmd *cobra.Command, args []string) {
		checkKeyFlags()
	},
	Long: `event runs a binary if there's a key change.`,
	Run:  startKey,
}

func startKey(cmd *cobra.Command, args []string) {
  Log("Just started the key watch.", "info")
  stdin := readStdin()
	if stdin != "null" {
    Log(fmt.Sprintf("stdin: '%s'", stdin), "info")
    Log("stdin wasn't blank", "info")
  } else {
    Log("stdin WAS blank", "info")
  }
}

func checkKeyFlags() {
	if Exec == "" {
		fmt.Println("Need a command to exec with '-e'")
		os.Exit(0)
	}
}

type ConsulKey struct {
	CreateIndex int    `json:"CreateIndex"`
	Flags       int    `json:"Flags"`
	Key         string `json:"Key"`
	LockIndex   int    `json:"LockIndex"`
	ModifyIndex int    `json:"ModifyIndex"`
	Session     string `json:"Session"`
	Value       string `json:"Value"`
}

func init() {
	RootCmd.AddCommand(keyCmd)
}
