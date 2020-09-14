package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// RootCmd is the root command
var RootCmd = &cobra.Command{
	Use:   "globe",
	Short: "Utility that glue together workflows",
	Long:  "Language-agnostic utility that glues configuration forutilities, task runners, and build tasks together",
}

// Execute adds all child commands to the root command and sets flags appropriately
func Execute() {
	if err := RootCmd.Execute(); err != nil {
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
