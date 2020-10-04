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
		// get data
		storeDir := cmd.Flag("dot-dir").Value.String()
		project := config.GetData(storeDir)
		userDotsConfig := config.GetUserToml(storeDir)

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

func resolveFile(src string, dest string, rel string) {
	fi, err := os.Lstat(dest)
	if err != nil {
		// file doesn't exist
		if os.IsNotExist(err) {
			err := os.Symlink(src, dest)
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
		logger.Debug("FILE has same content as thing. Replacing file with symlink")
		err := os.Remove(dest)
		util.P(err)

		err = os.Symlink(src, dest)
		util.P(err)
		return
	}

	// replace with link
	logger.Informational("FILE %s exists, is not a symlink and has different content. (remove, skip) ", src)
	var input string
	_, err = fmt.Scanln(&input)
	util.P(err)

	switch input {
	case "remove":
		err := os.Remove(dest)
		util.P(err)

		err = os.Symlink(src, dest)
		util.P(err)

	case "skip":
		logger.Debug("skipping '%s'\n", rel)
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
