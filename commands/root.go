package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

// RootCmd sets up all fo the other commands.
var RootCmd = &cobra.Command{
	Use:   "sifter",
	Short: "Make sure that a Consul watch only fires when it's appropriate.",
	Long:  `When Consul loads watches, it often fires them with a blank JSON array. sifter makes sure there's an actual JSON entity present, that it's new and then runs your command.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("`sifter -h` for help information.")
	},
}

var (
	// Exec is the command that should be executed if there's a watch that actually fires.
	Exec string

	// Token is used for access to Consul if an ACL is being used.
	Token string

	// DogStatsd sends metrics to Dogstatsd if set to true.
	DogStatsd bool

	// Prefix is the location in Consul's KV store to keep state information.
	Prefix string

	// MetricPrefix is the prefix for all Sifter metrics.
	MetricPrefix string

	// ConsulServer is where we want to connect to Consul.
	ConsulServer string
)

func init() {
	RootCmd.PersistentFlags().StringVarP(&Exec, "exec", "e", "", "Execute this command if a new event is present.")
	RootCmd.PersistentFlags().StringVarP(&Prefix, "prefix", "p", "sifter", "Consul prefix for saved state.")
	RootCmd.PersistentFlags().StringVarP(&Token, "token", "t", "anonymous", "Token for Consul access")
	RootCmd.PersistentFlags().StringVarP(&MetricPrefix, "metric", "m", "sifter", "Metric name for dogstatsd.")
	RootCmd.PersistentFlags().BoolVarP(&DogStatsd, "dogstatsd", "d", false, "send metrics to dogstatsd")
	RootCmd.PersistentFlags().StringVarP(&ConsulServer, "server", "s", "localhost:8500", "Consul server location")
}
