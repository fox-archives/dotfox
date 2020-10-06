package cmd

import (
	"github.com/eankeen/dotty/internal/util"
	"github.com/spf13/cobra"
)

var systemCmd = &cobra.Command{
	Use:   "system",
	Short: "Systemwide (/) config management",
	Long:  "Deals with system-wide (cross-user) configuration files. This may contain Refind, Systemwide XDG config, shell lists, etc.",
}

func init() {
	pf := systemCmd.PersistentFlags()
	pf.String("system-dir", "/", "Where to put dotfiles")

	err := cobra.MarkFlagDirname(pf, "system-dir")
	util.P(err)

	RootCmd.AddCommand(systemCmd)
}
