package cmd

import (
	"github.com/spf13/cobra"
)

var systemCmd = &cobra.Command{
	Use:   "system",
	Short: "Systemwide (/) config management",
	Long:  "Deals with system-wide (cross-user) configuration files. This may contain Refind, Systemwide XDG config, shell lists, etc.",
}

func init() {
	rootCmd.AddCommand(systemCmd)
}
