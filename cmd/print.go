package cmd

import (
	"fmt"

	"github.com/eankeen/dotty/config"
	"github.com/eankeen/dotty/internal/util"
	"github.com/spf13/cobra"
)

var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Print calculated variables",
	Long:  `Prints resolved locations of calculated variables`,
	Run: func(cmd *cobra.Command, args []string) {
		dotfilesDir := cmd.Flag("dotfiles-dir").Value.String()
		dottyCfg := config.DottyCfg(dotfilesDir)

		systemSrcDir := util.Src(dotfilesDir, dottyCfg, "system")
		systemDestDir := util.Dest(dotfilesDir, dottyCfg, "system")
		userSrcDir := util.Src(dotfilesDir, dottyCfg, "user")
		userDestDir := util.Dest(dotfilesDir, dottyCfg, "user")
		localSrcDir := util.Src(dotfilesDir, dottyCfg, "local")

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
