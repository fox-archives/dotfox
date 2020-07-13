package cmd

import (
	"github.com/eankeen/globe/config"
	"github.com/eankeen/globe/validate"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate to ensure your dotfiles configuration and files are valid",
	Long:  `Validate to ensure your dotfiles configuration and files are valid. This file also is ran before all other subcommands`,
	Run: func(cmd *cobra.Command, args []string) {
		// get data
		storeDir := cmd.Flag("store-dir").Value.String()
		projectDir := config.GetProjectLocation()
		project := config.GetData(projectDir, storeDir)

		// validate
		validate.Validate(validate.ValidationValues{
			StoreDir: storeDir,
			Project:  project,
		})
	},
}

func init() {
	RootCmd.AddCommand(checkCmd)
}
