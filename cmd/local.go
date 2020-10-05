package cmd

import (
	"github.com/spf13/cobra"
)

var localCmd = &cobra.Command{
	Use:   "local",
	Short: "Local (.) (per-project) config management",
	Long:  `Deal with configuration files contained in an independent project. This may contain EditorConfig, ESLint, Clang-Tidy etc. files`,
}

func init() {
	RootCmd.AddCommand(localCmd)
}
