package cmd

import (
	"fmt"
	"os"

	"github.com/eankeen/globe/scan"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd is the root command
var RootCmd = &cobra.Command{
	Use:   "globe",
	Short: "utility that glues",
	Long:  "Language-agnostic utility that glues configuration utilities, task runner, and build tasks together",
}

// Execute adds all child commands to the root command and sets flags appropriately
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(doInit)

	// RootCmd.PersistentFlags().StringVar("foo", "log-level", "", "Level for logging (info, warning (default), error")
}

func doInit() {
	viper.SetConfigFile(scan.GetConfigLocation())

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic("config file not found")
		}
		panic("some error occured")
	}

	fmt.Println("Using config file:", viper.ConfigFileUsed())
}
