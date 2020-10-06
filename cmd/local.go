package cmd

import (
	"os"

	"github.com/eankeen/dotty/internal/util"
	"github.com/spf13/cobra"
)

var localCmd = &cobra.Command{
	Use:   "local",
	Short: "Local (.) (per-project) config management",
	Long:  `Deal with configuration files contained in an independent project. This may contain EditorConfig, ESLint, Clang-Tidy etc. files`,
}

func init() {
	wd, err := os.Getwd()
	util.P(err)

	pf := systemCmd.PersistentFlags()
	pf.String("local-dir", wd, "Where to put dotfiles")

	err = cobra.MarkFlagDirname(pf, "local-dir")
	util.P(err)

	RootCmd.AddCommand(localCmd)
}
