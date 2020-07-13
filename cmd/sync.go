package cmd

import (
	"github.com/eankeen/globe/internal/util"
	"github.com/eankeen/globe/scan"
	"github.com/eankeen/globe/sync"
	"github.com/spf13/cobra"
)

var syncCommand = &cobra.Command{
	Use:   "sync",
	Short: "Sync Globe's configuration files",
	Long:  `Syncs configuration files`,
	Run: func(cmd *cobra.Command, args []string) {
		project := scan.Scan()
		util.PrintInfo("Project located at %s\n", project.ProjectLocation)
		sync.Sync(project)
	},
}

func init() {
	RootCmd.AddCommand(syncCommand)
}
