package cmd

import (
	"fmt"

	"github.com/eankeen/dotty/config"
	"github.com/spf13/cobra"
)

var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Print calculated variables",
	Long:  `Prints resolved locations of calculated variables`,
	Run: func(cmd *cobra.Command, args []string) {
		dotfilesDir := cmd.Flag("dotfiles-dir").Value.String()
		dottyCfg := config.DottyCfg(dotfilesDir)

		systemSrcDir := config.Src(dotfilesDir, dottyCfg, "system")
		systemDestDir := config.Dest(dotfilesDir, dottyCfg, "system")
		userSrcDir := config.Src(dotfilesDir, dottyCfg, "user")
		userDestDir := config.Dest(dotfilesDir, dottyCfg, "user")
		localSrcDir := config.Src(dotfilesDir, dottyCfg, "local")

		fmt.Printf("dotfilesDir: '%s'\n", dotfilesDir)
		fmt.Printf("systemSrcDir: '%s'\n", systemSrcDir)
		fmt.Printf("systemDestDir: '%s'\n", systemDestDir)
		fmt.Printf("userSrcDir: '%s'\n", userSrcDir)
		fmt.Printf("userDestDir: '%s'\n", userDestDir)
		fmt.Printf("localSrcDir: '%s'\n", localSrcDir)

	},
}

func init() {
	rootCmd.AddCommand(printCmd)
}
