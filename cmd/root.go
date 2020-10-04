package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// RootCmd is the root command
var RootCmd = &cobra.Command{
	Use:   "globe",
	Short: "Dotfile CM Utility",
	Long:  "A CM (Configuration Management) utility for dotfiles. Used for managing local, user, or system-wide dotfiles",
}

// Execute adds all child commands to the root command and sets flags appropriately
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	pf := RootCmd.PersistentFlags()
	pf.String("dot-dir", "", "The location of your dotfiles")
	err := cobra.MarkFlagDirname(pf, "dot-dir")
	if err != nil {
		panic(err)
	}
	err = cobra.MarkFlagRequired(pf, "dot-dir")
	if err != nil {
		panic(err)
	}
}
