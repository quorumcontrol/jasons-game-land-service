package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cmdService = &cobra.Command{
	Use:   "service",
	Short: "Service",
	Long:  `Service.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Service running\n")
	},
}
