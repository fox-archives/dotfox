package cmd

import (
	"fmt"
	"go-logger"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/eankeen/globe/config"
	"github.com/otiai10/copy"
	"github.com/spf13/cobra"
)

func resolveFile(src string, dest string, rel string) {
	fi, err := os.Lstat(dest)
	if err != nil {
		// file doesn't exist
		if os.IsNotExist(err) {
			err := os.Symlink(src, dest)
			if err != nil {
				panic(err)
			}
			return
		}
		panic(err)
	}

	// file exists and is a symbolic link
	if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
		linkDest, err := os.Readlink(dest)
		if err != nil {
			panic(err)
		}

		logger.PrintDebug("linkDest: %s\n", linkDest)
		logger.PrintDebug("src: %s\n", src)
		// symlink points to proper file
		if linkDest != src {
			logger.PrintDebug("Symlink does not match. Removing and recreating symlink\n")
			err := os.Remove(dest)
			if err != nil {
				panic(err)
			}

			err = os.Symlink(src, dest)
			if err != nil {
				panic(err)
			}

			return
		}

		// symlink does point to proper file,
		// nothing to do here
		return
	}

	// file exists and is NOT a symbolic link

	// read dest/src files to determine if they have the same content
	destContents, err := ioutil.ReadFile(dest)
	if err != nil {
		panic(err)
	}

	srcContents, err := ioutil.ReadFile(src)
	if err != nil {
		panic(err)
	}

	// if files have the same content
	if strings.Compare(string(destContents), string(srcContents)) == 0 {
		logger.PrintDebug("FILE has same content as thing. Replacing file with symlink")
		if err := os.Remove(dest); err != nil {
			panic(err)
		}

		if err := os.Symlink(src, dest); err != nil {
			panic(err)
		}
		return
	}

	// replace with link
	logger.PrintInfo("FILE %s exists, is not a symlink and has different content. (remove, skip) ", src)
	var input string
	_, err = fmt.Scanln(&input)
	if err != nil {
		panic(err)
	}

	switch input {
	case "remove":
		if err := os.Remove(dest); err != nil {
			panic(err)
		}

		if err := os.Symlink(src, dest); err != nil {
			panic(err)
		}

	case "skip":
		logger.PrintDebug("skipping '%s'\n", rel)
	}

	return
}

func resolveDirectory(src string, dest string, rel string) {
	fi, err := os.Lstat(dest)
	if err != nil {
		// file doesn't exist
		if os.IsNotExist(err) {
			// create folder symlink
			logger.PrintDebug("dest '%s' doesn't exist. recreating.\n", dest)

			err = os.Symlink(src, dest)
			if err != nil {
				panic(err)
			}
			return
		}

		panic(err)
	}

	// folder exists and is a symbolic link
	if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
		linkDest, err := os.Readlink(dest)
		if err != nil {
			panic(err)
		}

		// if link destination doesn't match src
		if linkDest != src {
			err := os.Remove(dest)
			if err != nil {
				panic(err)
			}

			err = os.Symlink(src, dest)
			if err != nil {
				panic(err)
			}
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
	if err != nil {
		panic(err)
	}

	err = os.Symlink(src, dest)
	if err != nil {
		panic(err)
	}

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
	if err != nil {
		panic(err)
	}

	var userDotsConfig UserDotsConfig
	if err := toml.Unmarshal(raw, &userDotsConfig); err != nil {
		panic(err)
	}

	return userDotsConfig
}

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Userwide (~) config management",
	Long:  "Actions to deal with configuration files that apply to a user's session. This may contain Bash startup, Vim config, X resource, etc. files",
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
					logger.PrintInfo("Operating on  File: '%s'\n", file.File)

					if info.IsDir() && file.Type == "folder" {
						resolveDirectory(src, dest, rel)
					} else if info.IsDir() && file.Type != "folder" {
						logger.PrintWarning("You expected '%s' (%s) to be a directory, but it's not", file.File, src)
					} else if !info.IsDir() && file.Type == "file" {
						resolveFile(src, dest, rel)
					} else if info.IsDir() && file.Type == "file" {
						logger.PrintWarning("'%s' is specified as a file, but at '%s', it is actually a directory\n", file.File, src)
					} else {
						logger.PrintWarning("Unexpected entry '%s' has type '%s' and isDir?: '%t'\n", file.File, file.Type, info.IsDir())
					}

					// we use first match
					return nil
				}
			}

			// match was not found
			// doesn't work because it false positives subdirs / files of folders
			// specified in configs, and random parent folders
			// logger.PrintWarning("File or folder '%s' was found in '%s', but it's not present in user.dots.toml\n", rel, dotfileDir)

			return nil
		})
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(userCmd)
}
