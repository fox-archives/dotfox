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
		dotDir := cmd.Flag("dot-dir").Value.String()
		srcDir := filepath.Join(dotDir, "user")
		destDir := cmd.Flag("user-dir").Value.String()

		actions.Apply(dotDir, srcDir, destDir)
	},
}

func init() {
	userCmd.AddCommand(userApplyCmd)
}
