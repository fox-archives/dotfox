package cmd

import (
	"fmt"
	"os"

	"github.com/eankeen/dotty/internal/util"
	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
)

// rootCmd is the root command
var rootCmd = &cobra.Command{
	Use:   "dotty",
	Short: "Dotfile CM Utility",
	Long:  "A CM (Configuration Management) utility for dotfiles. Used for managing local, user, or system-wide dotfiles",
}

// Execute adds all child commands to the root command and sets flags appropriately
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	unix.Umask(0022)

	pf := rootCmd.PersistentFlags()

	pf.String("dotfiles-dir", "", "The source locations of your dotfiles")
	cobra.MarkFlagRequired(pf, "dotfiles-dir")
	cobra.MarkFlagDirname(pf, "dotfiles-dir")

	pf.String("system-dir-dest", "/", "Destination of 'system' dotfiles")
	cobra.MarkFlagDirname(pf, "system-dir-dest")

	homedir, err := os.UserHomeDir()
	util.HandleError(err)
	pf.String("user-dir-dest", homedir, "Destination of 'user' dotfiles")
	cobra.MarkFlagDirname(pf, "user-dir-dest")

	wd, err := os.Getwd()
	util.HandleError(err)
	pf.String("local-dir-dest", wd, "Destination of 'local' dotfiles")
	cobra.MarkFlagDirname(pf, "local-dir-dest")
}
