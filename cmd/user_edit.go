package cmd

import (
	"github.com/eankeen/dotty/actions"
	"github.com/eankeen/dotty/config"
	"github.com/spf13/cobra"
)

var userEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit 'user' TOML config",
	Run: func(cmd *cobra.Command, args []string) {
		dotDir := cmd.Flag("dot-dir").Value.String()

		actions.OpenEditor(config.GetUserTomlPath(dotDir))
	},
}

func init() {
	userCmd.AddCommand(userEditCmd)
}
