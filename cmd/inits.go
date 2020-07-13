package cmd

import (
	"github.com/eankeen/globe/inits"
	"github.com/eankeen/globe/validate"
	"github.com/spf13/cobra"
)

var initsCmd = &cobra.Command{
	Use:   "init",
	Short: "Init Globe's configuration files",
	Long:  `Initiates configuration files to be used by Globe`,
	Run: func(cmd *cobra.Command, args []string) {
		v := validate.Validate(cmd, args)

		inits.Inits(v.StoreDir)
	},
}

func init() {
	RootCmd.AddCommand(initsCmd)
}
