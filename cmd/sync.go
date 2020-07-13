package cmd

import (
	"github.com/eankeen/globe/config"
	"github.com/eankeen/globe/sync"
	"github.com/eankeen/globe/validate"
	"github.com/spf13/cobra"
)

var syncCommand = &cobra.Command{
	Use:   "sync",
	Short: "Sync Globe's configuration files",
	Long:  `Syncs configuration files`,
	Run: func(cmd *cobra.Command, args []string) {
		// get data
		storeDir := cmd.Flag("store-dir").Value.String()
		projectDir := config.GetProjectLocation()
		project := config.GetData(projectDir, storeDir)

		// valudate values
		validate.Validate(validate.ValidationValues{
			StoreDir: storeDir,
			Project:  project,
		})

		// process files
		sync.ProcessFiles(project, project.SyncFiles.Files)
	},
}

func init() {
	RootCmd.AddCommand(syncCommand)
}
