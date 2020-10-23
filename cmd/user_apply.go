package cmd

import (
	"path/filepath"

	"github.com/eankeen/dotty/fs"
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

		onFile := func(src string, dest string, rel string) {
			fs.ApplyFile(src, dest, rel)
		}

		onFolder := func(src string, dest string, rel string) {
			fs.ApplyFolder(src, dest, rel)
		}

		fs.Walk(dotDir, srcDir, destDir, onFile, onFolder)
	},
}

func init() {
	userCmd.AddCommand(userApplyCmd)
}
