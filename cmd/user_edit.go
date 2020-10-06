package cmd

import (
	"github.com/eankeen/dotty/config"
	"github.com/eankeen/dotty/internal/util"
	"github.com/spf13/cobra"
)

var userEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit 'user' TOML config",
	Run: func(cmd *cobra.Command, args []string) {
		dotDir := cmd.Flag("dot-dir").Value.String()

		util.OpenEditor(config.GetUserTomlPath(dotDir))
	},
}

func init() {
	userCmd.AddCommand(userEditCmd)
}
