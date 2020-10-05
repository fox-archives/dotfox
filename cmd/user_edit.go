package cmd

import (
	"github.com/eankeen/globe/config"
	"github.com/eankeen/globe/internal/util"
	"github.com/spf13/cobra"
)

var userEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit 'user' TOML config",
	Run: func(cmd *cobra.Command, args []string) {
		storeDir := cmd.Flag("dot-dir").Value.String()

		util.OpenEditor(config.GetUserTomlPath(storeDir))
	},
}

func init() {
	userCmd.AddCommand(userEditCmd)
}