package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/eankeen/globe/config"
	"github.com/eankeen/globe/internal/util"
	logger "github.com/eankeen/go-logger"
	"github.com/spf13/cobra"
)

var userApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply updates intelligently",
	Run: func(cmd *cobra.Command, args []string) {
		// get data
		storeDir := cmd.Flag("dot-dir").Value.String()
		project := config.GetData(storeDir)
		userDotsConfig := readUserDotsConfig(project)

		dotfileDir := filepath.Join(project.StoreDir, "user")
		err := filepath.Walk(dotfileDir, func(path string, info os.FileInfo, err error) error {
			// prevent errors in slice
			if path == dotfileDir {
				return nil
			}

			src := path
			rel := path[len(dotfileDir)+1:]
			dest := filepath.Join(project.UserDir, rel)

			for _, file := range userDotsConfig.Files {
				if file.Type == "" {
					file.Type = "file"
				}

				if strings.HasSuffix(src, file.File) {
					logger.Informational("Operating on  File: '%s'\n", file.File)

					if info.IsDir() && file.Type == "folder" {
						resolveDirectory(src, dest, rel)
					} else if info.IsDir() && file.Type != "folder" {
						logger.Warning("You expected '%s' (%s) to be a directory, but it's not\n", file.File, src)
					} else if !info.IsDir() && file.Type == "file" {
						resolveFile(src, dest, rel)
					} else if info.IsDir() && file.Type == "file" {
						logger.Warning("'%s' is specified as a file, but at '%s', it is actually a directory\n", file.File, src)
					} else {
						logger.Warning("Unexpected entry '%s' has type '%s' and isDir?: '%t'\n", file.File, file.Type, info.IsDir())
					}

					// we use first match
					return nil
				}
			}

			// match was not found
			// doesn't work because it false positives subdirs / files of folders
			// specified in configs, and random parent folders
			// logger.Warning("File or folder '%s' was found in '%s', but it's not present in user.dots.toml\n", rel, dotfileDir)

			return nil
		})
		util.P(err)
	},
}

func init() {
	userCmd.AddCommand(userApplyCmd)
}
