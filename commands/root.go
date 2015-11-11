package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "sifter",
	Short: "Make sure there's a payload when Consul fires a watch.",
	Long:  `When Consul loads watches, it often fires them with a blank payload. sifter makes sure there's a payload present and runs your command.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("`sifter -h` for help information.")
	},
}

var (
	Exec string
)

func init() {
	RootCmd.PersistentFlags().StringVarP(&Exec, "exec", "e", "", "Execute this command if a payload is present.")
}
