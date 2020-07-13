package cmd

import (
	"github.com/eankeen/globe/config"
	"github.com/eankeen/globe/internal/util"
	"github.com/eankeen/globe/sync"
	"github.com/eankeen/globe/validate"
	"github.com/spf13/cobra"
)

var syncCommand = &cobra.Command{
	Use:   "sync",
	Short: "Sync Globe's configuration files",
	Long:  `Syncs configuration files`,
	Run: func(cmd *cobra.Command, args []string) {
		storeDir := cmd.Flag("store-dir").Value.String()
		projectDir := config.GetProjectLocation()

		validate.Validate(validate.ValidationValues{
			StoreDir: storeDir,
		})

		project := config.GetData(projectDir, storeDir)

		util.PrintInfo("Project located at %s\n", project.ProjectLocation)

		sync.ProcessFiles(project, project.SyncFiles.Files)
	},
}

func init() {
	RootCmd.AddCommand(syncCommand)
}
