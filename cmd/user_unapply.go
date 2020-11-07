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
		dotDir := cmd.Flag("dot-dir").Value.String()
		srcDir := filepath.Join(dotDir, "user")
		destDir := cmd.Flag("user-dir").Value.String()

		actions.Unapply(dotDir, srcDir, destDir)
	},
}

func init() {
	userCmd.AddCommand(userUnapplyCmd)
}
