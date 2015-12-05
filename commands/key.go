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
		watchEvent := decodeKeyStdin(stdin)
		watchEvent.examine()
		// What key are we watching?
		watchKey := watchEvent.getKey()
		// Grab the URL we will use to check Consul's previous SHA.
		nodeURL := watchEvent.makeURL()
		// Create the SHA256 from the value as passed on stdin.
		watchSHA := watchEvent.makeSHA()
		// Connect to Consul.
		consul, _ := Connect()
		// Get the previous value from Consul.
		previousSHA := Get(consul, nodeURL)
		if watchSHA != previousSHA {
			Set(consul, nodeURL, watchSHA)
			runCommand(Exec, "")
			RunTime(start, "complete", fmt.Sprintf("watch='key' exec='%s' sha='%s'", Exec, watchSHA))
			StatsdRunTime(start, Exec, "key", watchKey, watchSHA)
		} else {
			RunTime(start, "duplicate", fmt.Sprintf("watch='key' exec='%s' sha='%s'", Exec, watchSHA))
			StatsdDuplicate("key", watchKey)
		}
	} else {
		RunTime(start, "blank", fmt.Sprintf("watch='key' exec='%s'", Exec))
		StatsdBlank("key")
	}
}

func checkKeyFlags() {
	if Exec == "" {
		fmt.Println("Need a command to exec with '-e'")
		os.Exit(0)
	}
}

// KeyWatch is the JSON structure of a Consul key watch.
type KeyWatch struct {
	CreateIndex int    `json:"CreateIndex"`
	Flags       int    `json:"Flags,omitempty"`
	Key         string `json:"Key"`
	LockIndex   int    `json:"LockIndex,omitempty"`
	ModifyIndex int    `json:"ModifyIndex"`
	Session     string `json:"Session,omitempty"`
	Value       string `json:"Value"`
}

func decodeKeyStdin(data string) *KeyWatch {
	var key KeyWatch
	err := json.Unmarshal([]byte(data), &key)
	if err != nil {
		Log(fmt.Sprintf("error: %s", err), "info")
	}
	return &key
}

func (w *KeyWatch) examine() {
	created := w.CreateIndex
	modified := w.ModifyIndex
	keyName := w.getKey()
	value, _ := base64.StdEncoding.DecodeString(w.Value)
	Log(fmt.Sprintf("examine modifyIndex='%d' createIndex='%d' key='%s' value='%s'", modified, created, keyName, value), "debug")
}

func (w *KeyWatch) makeURL() string {
	hostname := getHostname()
	keyName := w.Key
	url := fmt.Sprintf("%s/key/%s/%s", Prefix, keyName, hostname)
	return url
}

func (w *KeyWatch) makeSHA() string {
	value, _ := base64.StdEncoding.DecodeString(w.Value)
	shaValue := sha256.Sum256([]byte(value))
	sha := fmt.Sprintf("%x", shaValue)
	return sha
}

func (w *KeyWatch) getKey() string {
	keyName := w.Key
	return keyName
}

func init() {
	RootCmd.AddCommand(keyCmd)
}
