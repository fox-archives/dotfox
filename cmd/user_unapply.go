package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/eankeen/dotty/config"
	"github.com/eankeen/dotty/internal/util"
	logger "github.com/eankeen/go-logger"
	"github.com/spf13/cobra"
)

var userUnapplyCmd = &cobra.Command{
	Use:   "unapply",
	Short: "Unapply a",
	Long:  "This unapplies all user dotfiles, unlinking them from the destination (user) directory",
	Run: func(cmd *cobra.Command, args []string) {
		dotDir := cmd.Flag("dot-dir").Value.String()
		destDir := cmd.Flag("user-dir").Value.String()

		userDir := filepath.Join(dotDir, "user")
		userToml := config.GetUserToml(dotDir)

		err := filepath.Walk(userDir, func(src string, srcInfo os.FileInfo, err error) error {
			// prevent errors in slice
			if src == userDir {
				return nil
			}

			rel := src[len(userDir)+1:]
			dest := filepath.Join(destDir, rel)

			for _, file := range userToml.Files {
				logger.Debug("src: %s\n", src)
				logger.Debug("file.File: %s\n", file.File)

				// if path has a part in ignores, then we skip the whole file
				for _, ignore := range userToml.Ignores {
					ignoreEntryMatches, _ := config.IgnoreMatches(src, ignore)

					if ignoreEntryMatches {
						logger.Debug("Ignoring '%s'\n", src)
						return nil
					}
				}

				if len(file.File) == 0 {
					logger.Warning("An entry in your '%s' config does not have a 'name' property. (This message may be repeated multiple times for each entry). Skipping\n", "user.dots.toml")
					return nil
				}

				fileMatches, fileType := config.FileMatches(src, file)

				if fileMatches && fileType == "file" {
					logger.Informational("Operating on File: '%s'\n", file.File)

					if srcInfo.IsDir() {
						logger.Warning("Your '%s' entry has a match, but it actually is a folder (%s) instead of a file. Did you mean to append a slash? Skipping file", file.File, src)
						return nil
					}

					resolveFile(src, dest, rel)

					// go to next file (in dotfile folder)
					return nil
				} else if fileMatches && fileType == "folder" {
					logger.Informational("Unapplying on Folder: '%s'\n", file.File)

					if !srcInfo.IsDir() {
						logger.Warning("Your '%s' entry has a match, but it actually matches a file (%s) instead of a folder. Did you mean to remove the trailing slack? Skipping file\n", file.File, src)
						return nil
					}

					unapplyDirectory(src, dest, rel)

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
	userCmd.AddCommand(userUnapplyCmd)

}

func unapplyDirectory(src string, dest string, rel string) {
	fi, err := os.Lstat(dest)
	if os.IsNotExist(err) {
		return
	}
	util.HandleFsError(err)

	if fi.Mode()&os.ModeSymlink != os.ModeSymlink {
		logger.Error("Skipping: Not a Symlink: '%s'\n", dest)
		return
	}

	cmd := exec.Command("unlink", dest)
	res, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(res)
		util.HandleError(err)
	}

	// if fi.IsDir() {
	// 	return
	// }

	// realDest, err := os.Readlink(dest)
	// util.HandleFsError(err)

	// fmt.Println(realDest)
	// err = syscall.Unlink(realDest)
	// util.HandleError(err)
}
