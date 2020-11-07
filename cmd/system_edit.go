package cmd

import (
	"github.com/eankeen/dotty/actions"
	"github.com/eankeen/dotty/config"
	"github.com/spf13/cobra"
)

var systemEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit 'system' TOML config",
	Run: func(cmd *cobra.Command, args []string) {
		dotfilesDir := cmd.Flag("dotfiles-dir").Value.String()

		actions.OpenEditor(config.GetSystemTomlPath(dotfilesDir))
	},
}

func init() {
	systemCmd.AddCommand(systemEditCmd)
}
