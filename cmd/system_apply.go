package cmd

import (
	"os"
	"path/filepath"

	"github.com/eankeen/dotty/fs"
	"github.com/eankeen/go-logger"
	"github.com/spf13/cobra"
)

var systemApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Synchronize system dotfiles",
	Long:  "Synchronize system dotfiles. You will be prompted on conflicts",
	Run: func(cmd *cobra.Command, args []string) {
		if os.Getuid() != 0 {
			logger.Error("Must run as root. Exiting\n")
			os.Exit(1)
		}

		dotDir := cmd.Flag("dot-dir").Value.String()
		srcDir := filepath.Join(dotDir, "system")
		destDir := cmd.Flag("system-dir").Value.String()

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
	systemCmd.AddCommand(systemApplyCmd)
}
