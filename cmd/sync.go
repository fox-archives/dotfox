package cmd

import (
	"github.com/eankeen/globe/sync"
	"github.com/eankeen/globe/validate"
	"github.com/spf13/cobra"
)

var syncCommand = &cobra.Command{
	Use:   "sync",
	Short: "Sync Globe's configuration files",
	Long:  `Syncs configuration files`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		validate.Validate(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		sync.Sync()
	},
}

func init() {
	RootCmd.AddCommand(syncCommand)
}
