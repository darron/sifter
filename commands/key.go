package commands

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var keyCmd = &cobra.Command{
	Use:   "key",
	Short: "Run a binary from a Consul key watch.",
	PreRun: func(cmd *cobra.Command, args []string) {
		checkKeyFlags()
	},
	Long: `event runs a binary if there's a key change.`,
	Run:  startKey,
}

func startKey(cmd *cobra.Command, args []string) {
	start := time.Now()
	stdin := readStdin()
	if stdin != "null" {
		d := decodeKeyStdin(stdin)
		modified, key, value := d.examine()
		Log(fmt.Sprintf("mod='%d' key='%s' value='%s'", modified, key, value), "info")
	} else {
		RunTime(start, "blank", fmt.Sprintf("watch='key' exec='%s'", Exec))
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
	Flags       int    `json:"Flags,omitempty"`
	Key         string `json:"Key"`
	LockIndex   int    `json:"LockIndex,omitempty"`
	ModifyIndex int    `json:"ModifyIndex"`
	Session     string `json:"Session,omitempty"`
	Value       string `json:"Value"`
}

func decodeKeyStdin(data string) *ConsulKey {
	var key ConsulKey
	err := json.Unmarshal([]byte(data), &key)
	if err != nil {
		Log(fmt.Sprintf("error: %s", err), "info")
	}
	return &key
}

func (c *ConsulKey) examine() (int, string, string) {
	modified := c.ModifyIndex
	keyName := c.Key
	value, _ := base64.StdEncoding.DecodeString(c.Value)
	return modified, keyName, string(value)
}

func init() {
	RootCmd.AddCommand(keyCmd)
}
