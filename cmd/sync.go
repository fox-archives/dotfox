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
		validatedArgs := validate.Validate(storeDir)

		project := config.GetData(config.GetProjectLocation(), validatedArgs.StoreDir)

		util.PrintInfo("Project located at %s\n", project.ProjectLocation)

		for _, file := range project.SyncFiles.Files {
			util.PrintInfo("Processing file %s\n", file.RelPath)

			if file.Op == "add" {
				sync.CopyFile(project, file)
				continue
			} else if file.Op == "remove" {
				sync.RemoveFile(project, file)
				continue
			}

			util.PrintError("File '%s's operation could not be read. Exiting.\n", file.RelPath)
		}
	},
}

func init() {
	RootCmd.AddCommand(syncCommand)
}
