package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var systemApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply updates intelligently",
	Run: func(cmd *cobra.Command, args []string) {
		if os.Getuid() != 0 {
			fmt.Println("Please run with 'sudo' and try again")
			os.Exit(1)
		}
	},
}

func init() {
	systemCmd.AddCommand(systemApplyCmd)
}
