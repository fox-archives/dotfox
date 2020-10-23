package cmd

import (
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Userwide (~) config management",
	Long:  "Actions to deal with configuration files that apply to a user's session. This may contain Bash startup, Vim config, X resource, etc. files",
}

func init() {
	rootCmd.AddCommand(userCmd)
}
