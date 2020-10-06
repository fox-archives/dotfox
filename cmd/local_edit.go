package cmd

import (
	"github.com/eankeen/dotty/config"
	"github.com/eankeen/dotty/internal/util"
	"github.com/spf13/cobra"
)

var localEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit 'local' TOML config",
	Run: func(cmd *cobra.Command, args []string) {
		dotDir := cmd.Flag("dot-dir").Value.String()

		util.OpenEditor(config.GetLocalTomlPath(dotDir))
	},
}

func init() {
	localCmd.AddCommand(localEditCmd)
}
