package cmd

import (
	"path/filepath"

	"github.com/eankeen/dotty/fs"
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

		onFile := func(src string, dest string, rel string) {
			fs.UnapplyFile(src, dest, rel)
		}

		onFolder := func(src string, dest string, rel string) {
			fs.UnapplyFolder(src, dest, rel)
		}

		fs.Walk(dotDir, srcDir, destDir, onFile, onFolder)
	},
}

func init() {
	userCmd.AddCommand(userUnapplyCmd)
}
