package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/eankeen/globe/config"
	"github.com/eankeen/globe/internal/util"
	logger "github.com/eankeen/go-logger"
	"github.com/otiai10/copy"
	"github.com/spf13/cobra"
)

var userApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply updates intelligently",
	Run: func(cmd *cobra.Command, args []string) {
		storeDir := cmd.Flag("dot-dir").Value.String()
		destDir := cmd.Flag("dest-dir").Value.String()
		userDir := filepath.Join(storeDir, "user")

		userDotsConfig := config.GetUserToml(storeDir)

		err := filepath.Walk(userDir, func(src string, srcInfo os.FileInfo, err error) error {
			// prevent errors in slice
			if src == userDir {
				return nil
			}

			rel := src[len(userDir)+1:]
			dest := filepath.Join(destDir, rel)

			for _, file := range userDotsConfig.Files {
				fileMatches, fileType := config.FileMatches(src, srcInfo, file)

				if fileMatches && fileType == "file" {
					logger.Informational("Operating on  File: '%s'\n", file.File)
					resolveFile(src, dest, rel)

					// go to next file (in dotfile folder)
					return nil
				} else if fileMatches && fileType == "folder" {
					logger.Informational("Operating on  Folder: '%s'\n", file.File)
					resolveDirectory(src, dest, rel)

					// go to next file (in dotfile folder)
					return nil
				}

				// file in config did not match. continue because
				// another one might
				continue
			}

			// no explicit match was made. it may have already been matched, not match at all, or a parent folder matched instead
			return nil
		})
		util.P(err)
	},
}

func init() {
	userCmd.AddCommand(userApplyCmd)
}

func resolveFile(src string, dest string, rel string) {
	fi, err := os.Lstat(dest)
	if err != nil {
		// file doesn't exist
		if os.IsNotExist(err) {
			err := os.MkdirAll(filepath.Dir(dest), 0755)
			util.P(err)

			err = os.Symlink(src, dest)
			util.P(err)
			return
		}
		panic(err)
	}

	// file exists and is a symbolic link
	if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
		linkDest, err := os.Readlink(dest)
		util.P(err)

		logger.Debug("linkDest: %s\n", linkDest)
		logger.Debug("src: %s\n", src)
		// symlink points to proper file
		if linkDest != src {
			logger.Debug("Symlink does not match. Removing and recreating symlink\n")
			err := os.Remove(dest)
			util.P(err)

			err = os.Symlink(src, dest)
			util.P(err)

			return
		}

		// symlink does point to proper file,
		// nothing to do here
		return
	}

	// file exists and is NOT a symbolic link

	// read dest/src files to determine if they have the same content
	destContents, err := ioutil.ReadFile(dest)
	util.P(err)

	srcContents, err := ioutil.ReadFile(src)
	util.P(err)

	// if files have the same content
	if strings.Compare(string(destContents), string(srcContents)) == 0 {
		logger.Debug("FILE has same content as thing. Replacing file with symlink\n")
		err := os.Remove(dest)
		util.P(err)

		err = os.Symlink(src, dest)
		util.P(err)
		return
	}

	// replace with link
	logger.Informational("FILE %s exists, is not a symlink and has different content. (overwrite, skip,overwrite) ", src)
	var input string
	_, err = fmt.Scanln(&input)
	util.P(err)

	switch input {
	case "overwrite":
		err := os.Remove(dest)
		util.P(err)

		err = os.Symlink(src, dest)
		util.P(err)
		break
	case "skip":
		logger.Debug("skipping '%s'\n", rel)
		break
	default:
		logger.Informational("Unknown Response\n")
	}

	return
}

func resolveDirectory(src string, dest string, rel string) {
	fi, err := os.Lstat(dest)
	if err != nil {
		// file doesn't exist
		if os.IsNotExist(err) {
			// create folder symlink
			logger.Debug("dest '%s' doesn't exist. recreating.\n", dest)

			err = os.Symlink(src, dest)
			util.P(err)

			return
		}

		panic(err)
	}

	// folder exists and is a symbolic link
	if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
		linkDest, err := os.Readlink(dest)
		util.P(err)

		// if link destination doesn't match src
		if linkDest != src {
			err := os.Remove(dest)
			util.P(err)

			err = os.Symlink(src, dest)
			util.P(err)
		}

		// link has correct dest, no need to do anything
		return
	}

	// folder exists and is not a symbolic link
	err = copy.Copy(dest, src)
	if err != nil {
		// if file already exists, we don't overwrite
		// this might want to be changed later
		if os.IsExist(err) {

		} else {
			panic(err)
		}
	}

	err = os.RemoveAll(dest)
	util.P(err)

	err = os.Symlink(src, dest)
	util.P(err)

	return
}
