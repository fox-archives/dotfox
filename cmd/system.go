package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var systemCmd = &cobra.Command{
	Use:   "init",
	Short: "Systemwide (/) config management",
	Long:  "Deals with system-wide (cross-user) configuration files. This may contain Refind, Systemwide XDG config, shell lists, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		if os.Getuid() != 0 {
			fmt.Println("Please run with 'sudo' and try again")
			os.Exit(1)
		}
	},
}
