package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/eankeen/dotty/config"
	"github.com/eankeen/dotty/internal/util"
	"github.com/eankeen/go-logger"
	"github.com/spf13/cobra"
)

var systemApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply updates intelligently",
	Run: func(cmd *cobra.Command, args []string) {
		if os.Getuid() != 0 {
			logger.Error("Must run as root. Exiting\n")
			os.Exit(1)
		}

		dotDir := cmd.Flag("dot-dir").Value.String()
		systemToml := config.GetSystemToml(dotDir)

		for _, file := range systemToml.Files {
			src := filepath.Join(dotDir, "system", file.File)
			dst := file.File

			srcFi, err := os.Stat(src)
			util.HandleFsError(err)

			if srcFi.IsDir() {
				logger.Informational("Skipping directory %s\n", file.File)
				continue
			}

			// if is file
			logger.Informational("Processing %s\n", file.File)
			srcContents, err := ioutil.ReadFile(src)
			util.HandleFsError(err)

			err = ioutil.WriteFile(dst, srcContents, 0644)
			util.HandleFsError(err)
		}

	},
}

func init() {
	systemCmd.AddCommand(systemApplyCmd)
}
