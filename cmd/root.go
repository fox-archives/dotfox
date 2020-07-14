package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/eankeen/globe/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	cobra.OnInitialize(func() {
		projectDir := config.GetProjectDir()
		globeConfig := path.Join(projectDir, "globe.toml")

		viper.SetConfigFile(globeConfig)
	})

	pf := RootCmd.PersistentFlags()
	pf.String("store-dir", "", "The location of your dotfiles")
	err := cobra.MarkFlagDirname(pf, "store-dir")
	if err != nil {
		panic(err)
	}
	err = cobra.MarkFlagRequired(pf, "store-dir")
	if err != nil {
		panic(err)
	}

	// RootCmd.PersistentFlags().StringVar("foo", "log-level", "", "Level for logging (info, warning (default), error")
}
