package commands

import (
	"crypto/sha256"
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
		d.examine()
		// Grab the URL we will use to check Consul.
		url := d.getUrl()
		// Create the SHA256 from the value as passed on stdin.
		shaValue := d.getSHA()
		// Get the previous value from Consul.
		c, _ := Connect()
		urlData := Get(c, url)
		if shaValue != urlData {
			Set(c, url, shaValue)
			runCommand(Exec, "")
			RunTime(start, "complete", fmt.Sprintf("watch='key' exec='%s' sha='%s'", Exec, shaValue))
			StatsdRunTime(start, Exec, "key", d.getKey(), shaValue)
		} else {
			RunTime(start, "duplicate", fmt.Sprintf("watch='key' exec='%s' sha='%s'", Exec, shaValue))
		}
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

func (c *ConsulKey) examine() {
	modified := c.ModifyIndex
	keyName := c.Key
	value, _ := base64.StdEncoding.DecodeString(c.Value)
	Log(fmt.Sprintf("mod='%d' key='%s' value='%s'", modified, keyName, value), "debug")
}

func (c *ConsulKey) getUrl() string {
	hostname := getHostname()
	keyName := c.Key
	url := fmt.Sprintf("%s/key/%s/%s", Prefix, keyName, hostname)
	return url
}

func (c *ConsulKey) getSHA() string {
	value, _ := base64.StdEncoding.DecodeString(c.Value)
	shaValue := sha256.Sum256([]byte(value))
	sha := fmt.Sprintf("%x", shaValue)
	return sha
}

func (c *ConsulKey) getKey() string {
	keyName := c.Key
	return keyName
}

func init() {
	RootCmd.AddCommand(keyCmd)
}
