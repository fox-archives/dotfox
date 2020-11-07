package cmd

import (
	"github.com/eankeen/dotty/actions"
	"github.com/eankeen/dotty/config"
	"github.com/spf13/cobra"
)

var localEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit 'local' TOML config",
	Run: func(cmd *cobra.Command, args []string) {
		dotDir := cmd.Flag("dot-dir").Value.String()

		actions.OpenEditor(config.GetLocalTomlPath(dotDir))
	},
}

func init() {
	localCmd.AddCommand(localEditCmd)
}
