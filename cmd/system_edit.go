package cmd

import (
	"github.com/eankeen/globe/config"
	"github.com/eankeen/globe/internal/util"
	"github.com/spf13/cobra"
)

var systemEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit 'system' TOML config",
	Run: func(cmd *cobra.Command, args []string) {
		storeDir := cmd.Flag("dot-dir").Value.String()

		util.OpenEditor(config.GetSystemTomlPath(storeDir))
	},
}

func init() {
	systemCmd.AddCommand(systemEditCmd)
}
