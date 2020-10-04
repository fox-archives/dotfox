package cmd

import (
	"github.com/eankeen/globe/config"
	"github.com/spf13/cobra"
)

var localCmd = &cobra.Command{
	Use:   "local",
	Short: "Local (.) (per-project) config management",
	Long:  `Deal with configuration files contained in an independent project. This may contain EditorConfig, ESLint, Clang-Tidy etc. files`,
	Run: func(cmd *cobra.Command, args []string) {
		// write globe.state
		writeGlobeState()

		// get data
		storeDir := cmd.Flag("dot-dir").Value.String()
		project := config.GetData(storeDir)

		// process filesproject
		ProcessFiles(project, project.Files)
	},
}

func init() {
	RootCmd.AddCommand(localCmd)
}
