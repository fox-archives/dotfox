package cmd

import (
	"path/filepath"

	"github.com/eankeen/dotty/actions"
	"github.com/eankeen/dotty/config"
	"github.com/spf13/cobra"
)

var userEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit 'user' TOML config",
	Run: func(cmd *cobra.Command, args []string) {
		dotfilesDir := cmd.Flag("dotfiles-dir").Value.String()

		file := filepath.Join(dotfilesDir, config.DottyCfg(dotfilesDir).ConfigDir, "user.dots.toml")
		actions.OpenEditor(file)
	},
}

func init() {
	userCmd.AddCommand(userEditCmd)
}
