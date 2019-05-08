package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "land-service",
	Short:        "Jason's Game land service",
	Long:         `Jason's Game land service`,
	SilenceUsage: true,
}
var configFilePath string

func Execute() {
	rootCmd.Flags().StringVarP(&configFilePath, "config", "c", "conf.json",
		"Configuration file")
	rootCmd.AddCommand(cmdService, cmdClient, cmdGenerateKey, cmdBootstrapper)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
