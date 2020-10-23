package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/eankeen/dotty/config"
	"github.com/eankeen/dotty/fs"
	"github.com/eankeen/dotty/internal/util"
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
		destDir := cmd.Flag("system-dir").Value.String()

		systemDir := filepath.Join(dotDir, "system")
		systemToml := config.GetSystemToml(dotDir)

		err := filepath.Walk(systemDir, func(src string, srcInfo os.FileInfo, err error) error {
			// prevent errors in slice
			if src == systemDir {
				return nil
			}

			rel := src[len(systemDir)+1:]
			dest := filepath.Join(destDir, rel)

			fmt.Println(rel)

			for _, file := range systemToml.Files {
				logger.Debug("src: %s\n", src)
				logger.Debug("file.File: %s\n", file.File)

				// if path has a part in ignores, then we skip the whole file
				for _, ignore := range systemToml.Ignores {
					ignoreEntryMatches, _ := config.IgnoreMatches(src, ignore)

					if ignoreEntryMatches {
						logger.Debug("Ignoring '%s'\n", src)
						return nil
					}
				}

				if len(file.File) == 0 {
					logger.Warning("An entry in your '%s' config does not have a 'name' property. (This message may be repeated multiple times for each entry). Skipping\n", "system.dots.toml")
					return nil
				}

				fileMatches, fileType := config.FileMatches(src, file)

				if fileMatches && fileType == "file" {
					logger.Informational("Operating on File: '%s'\n", file.File)

					if srcInfo.IsDir() {
						logger.Warning("Your '%s' entry has a match, but it actually is a folder (%s) instead of a file. Did you mean to append a slash? Skipping file", file.File, src)
						return nil
					}

					fs.ApplyFile(src, dest, rel)

					// go to next file (in dotfile folder)
					return nil
				} else if fileMatches && fileType == "folder" {
					logger.Informational("Operating on Folder: '%s'\n", file.File)

					if !srcInfo.IsDir() {
						logger.Warning("Your '%s' entry has a match, but it actually matches a file (%s) instead of a folder. Did you mean to remove the trailing slack? Skipping file\n", file.File, src)
						return nil
					}

					fs.ApplyDirectory(src, dest, rel)

					// go to next file (in dotfile folder)
					return nil
				}

				// file in config did not match. continue because
				// another one might
			}

			// no explicit match was made. it may have already been matched, not match at all, or a parent folder matched instead
			return nil
		})
		util.HandleFsError(err)

	},
}

func init() {
	systemCmd.AddCommand(systemApplyCmd)
}
