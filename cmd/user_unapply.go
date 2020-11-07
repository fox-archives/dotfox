package cmd

import (
	"path/filepath"

	"github.com/eankeen/dotty/actions"
	"github.com/spf13/cobra"
)

var userUnapplyCmd = &cobra.Command{
	Use:   "unapply",
	Short: "Unapply a",
	Long:  "This unapplies all user dotfiles, unlinking them from the destination (user) directory",
	Run: func(cmd *cobra.Command, args []string) {
		dotfilesDir := cmd.Flag("dotfiles-dir").Value.String()
		srcDir := filepath.Join(dotfilesDir, "user")
		destDir := cmd.Flag("user-dir").Value.String()

		actions.Unapply(dotfilesDir, srcDir, destDir)
	},
}

func init() {
	userCmd.AddCommand(userUnapplyCmd)
}
