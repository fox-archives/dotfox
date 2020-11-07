package cmd

import (
	"path/filepath"

	"github.com/eankeen/dotty/actions"
	"github.com/spf13/cobra"
)

var userApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Synchronize user dotfiles",
	Long:  "Synchronize user dotfiles. You will be prompted on conflicts",
	Run: func(cmd *cobra.Command, args []string) {
		dotfilesDir := cmd.Flag("dotfiles-dir").Value.String()
		srcDir := filepath.Join(dotfilesDir, "user")
		destDir := cmd.Flag("user-dir").Value.String()

		actions.Apply(dotfilesDir, srcDir, destDir)
	},
}

func init() {
	userCmd.AddCommand(userApplyCmd)
}
