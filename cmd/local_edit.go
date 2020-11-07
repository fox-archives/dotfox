package cmd

import (
	"path/filepath"

	"github.com/eankeen/dotty/actions"
	"github.com/eankeen/dotty/config"
	"github.com/spf13/cobra"
)

var localEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit 'local' TOML config",
	Run: func(cmd *cobra.Command, args []string) {
		dotfilesDir := cmd.Flag("dotfiles-dir").Value.String()

		file := filepath.Join(dotfilesDir, config.DottyCfg(dotfilesDir).ConfigDir, "local.dots.toml")
		actions.OpenEditor(file)
	},
}

func init() {
	localCmd.AddCommand(localEditCmd)
}
