package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "sifter",
	Short: "Make sure that a Consul watch only fires when it's appropriate.",
	Long:  `When Consul loads watches, it often fires them with a blank JSON array. sifter makes sure there's an actual JSON entity present, that it's new and then runs your command.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("`sifter -h` for help information.")
	},
}

var (
	Exec      string
	Token     string
	DogStatsd bool
	Prefix    string
)

func init() {
	RootCmd.PersistentFlags().StringVarP(&Exec, "exec", "e", "", "Execute this command if a new event is present.")
	RootCmd.PersistentFlags().StringVarP(&Prefix, "prefix", "p", "sifter", "Consul prefix for saved state.")
	RootCmd.PersistentFlags().StringVarP(&Token, "token", "t", "anonymous", "Token for Consul access")
	RootCmd.PersistentFlags().BoolVarP(&DogStatsd, "dogstatsd", "d", false, "send metrics to dogstatsd")
}
