package cmd

import (
	"github.com/eankeen/globe/validate"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate to ensure your dotfiles configuration and files are valid",
	Long:  `Validate to ensure your dotfiles configuration and files are valid. This file also is ran before all other subcommands`,
	Run: func(cmd *cobra.Command, args []string) {
		validate.Validate(cmd, args)
	},
}

func init() {
	RootCmd.AddCommand(checkCmd)
}
