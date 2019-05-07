package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "land-service",
	Short: "Jason's Game land service",
	Long:  `Jason's Game land service`,
}

func Execute() {
	rootCmd.AddCommand(cmdService, cmdClient)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
