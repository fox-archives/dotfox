package cmd

import (
	"github.com/eankeen/globe/config"
	"github.com/eankeen/globe/internal/util"
	"github.com/spf13/cobra"
)

var localEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit 'local' TOML config",
	Run: func(cmd *cobra.Command, args []string) {
		storeDir := cmd.Flag("dot-dir").Value.String()

		util.OpenEditor(config.GetLocalTomlPath(storeDir))
	},
}

func init() {
	localCmd.AddCommand(localEditCmd)
}
