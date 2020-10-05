package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/eankeen/globe/config"
	"github.com/eankeen/globe/internal/util"
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

		storeDir := cmd.Flag("dot-dir").Value.String()
		systemToml := config.GetSystemToml(storeDir)

		for _, file := range systemToml.Files {
			src := filepath.Join(storeDir, "system", file.File)
			dst := file.File

			srcFi, err := os.Stat(src)
			util.P(err)

			if srcFi.IsDir() {
				logger.Informational("Skipping directory %s\n", file.File)
				continue
			}

			// if is file
			logger.Informational("Processing %s\n", file.File)
			srcContents, err := ioutil.ReadFile(src)
			util.P(err)

			err = ioutil.WriteFile(dst, srcContents, 0644)
			util.P(err)
		}

	},
}

func init() {
	systemCmd.AddCommand(systemApplyCmd)
}
