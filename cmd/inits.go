package cmd

import (
	"github.com/eankeen/globe/inits"
	"github.com/eankeen/globe/validate"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init Globe's configuration files",
	Long:  `Initiates configuration files to be used by Globe`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		validate.Validate(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		inits.Inits()
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}
