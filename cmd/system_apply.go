package cmd

import (
	"os"
	"path/filepath"

	"github.com/eankeen/dotty/actions"
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

		actions.Apply(dotDir, srcDir, destDir)
	},
}

func init() {
	systemCmd.AddCommand(systemApplyCmd)
}
