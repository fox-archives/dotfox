package cmd

import (
	"github.com/eankeen/dotty/actions"
	"github.com/eankeen/dotty/config"
	"github.com/eankeen/dotty/internal/util"
	"github.com/spf13/cobra"
)

var userApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Synchronize user dotfiles",
	Long:  "Synchronize user dotfiles. You will be prompted on conflicts",
	Run: func(cmd *cobra.Command, args []string) {
		dotfilesDir := cmd.Flag("dotfiles-dir").Value.String()
		dottyCfg := config.DottyCfg(dotfilesDir)

		srcDir := util.Src(dotfilesDir, dottyCfg, "user")
		destDir := util.Dest(dotfilesDir, dottyCfg, "user")

		actions.Apply(dotfilesDir, srcDir, destDir)
	},
}

func init() {
	userCmd.AddCommand(userApplyCmd)
}
