package cmd

import (
	"github.com/spf13/cobra"
)

var cmdClient = &cobra.Command{
	Use:   "client",
	Short: "Jason's Game land service client",
}

var clientPort uint

func init() {
	rootCmd.AddCommand(cmdClient)
	cmdClient.PersistentFlags().UintVarP(&clientPort, "port", "p", 0, "Source port")
}
