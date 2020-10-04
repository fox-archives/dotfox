package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/eankeen/globe/config"
	"github.com/eankeen/globe/internal/util"
	logger "github.com/eankeen/go-logger"
	"github.com/otiai10/copy"
	"github.com/spf13/cobra"
)

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

type File struct {
	File string   `toml:"file"`
	Tags []string `toml:"tags"`
	Type string   `toml:"type"`
}

type UserDotsConfig struct {
	Files []File `toml:"files"`
}

func readUserDotsConfig(project config.Project) UserDotsConfig {
	projectConfig := filepath.Join(project.StoreDir, "user.dots.toml")

	raw, err := ioutil.ReadFile(projectConfig)
	util.P(err)

	var userDotsConfig UserDotsConfig
	err = toml.Unmarshal(raw, &userDotsConfig)
	util.P(err)

	return userDotsConfig
}

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Userwide (~) config management",
	Long:  "Actions to deal with configuration files that apply to a user's session. This may contain Bash startup, Vim config, X resource, etc. files",
}

func init() {
	RootCmd.AddCommand(userCmd)
}
