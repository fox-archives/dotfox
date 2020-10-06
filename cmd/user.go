package cmd

import (
	"os"

	"github.com/eankeen/dotty/internal/util"
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Userwide (~) config management",
	Long:  "Actions to deal with configuration files that apply to a user's session. This may contain Bash startup, Vim config, X resource, etc. files",
}

func init() {
	homedir, err := os.UserHomeDir()
	util.P(err)

	pf := systemCmd.PersistentFlags()
	pf.String("user-dir", homedir, "Where to put dotfiles")

	err = cobra.MarkFlagDirname(pf, "user-dir")
	util.P(err)

	RootCmd.AddCommand(userCmd)
}
