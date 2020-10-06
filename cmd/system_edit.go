package cmd

import (
	"github.com/eankeen/dotty/config"
	"github.com/eankeen/dotty/internal/util"
	"github.com/spf13/cobra"
)

var systemEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit 'system' TOML config",
	Run: func(cmd *cobra.Command, args []string) {
		dotDir := cmd.Flag("dot-dir").Value.String()

		util.OpenEditor(config.GetSystemTomlPath(dotDir))
	},
}

func init() {
	systemCmd.AddCommand(systemEditCmd)
}
