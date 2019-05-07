package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cmdClient = &cobra.Command{
	Use:   "client",
	Short: "Jason's Game land service client",
	Long:  "Jason's Game land service client",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Client running\n")
	},
}
