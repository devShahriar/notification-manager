package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "notification-manager",
	Short: "notification-manager",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Print(err)
	}
}
