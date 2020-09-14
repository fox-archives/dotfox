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
	"github.com/spf13/cobra"
)

func resolveFile(src string, dest string, rel string) error {
	fi, err := os.Lstat(dest)
	if err != nil {
		// file doesn't exist
		if os.IsNotExist(err) {
			if err := os.Symlink(src, dest); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	// file exists and is a symbolic link
	if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
		linkDest, err := os.Readlink(dest)
		if err != nil {
			return err
		}

		logger.PrintDebug("linkDest: %s\n", linkDest)
		logger.PrintDebug("src: %s\n", src)
		// symlink points to proper file
		if linkDest != src {
			logger.PrintDebug("Symlink does not match. Removing and recreating symlink\n")
			if err := os.Remove(dest); err != nil {
				return err
			}

			if err := os.Symlink(src, dest); err != nil {
				return err
			}

			return nil
		}

		// symlink does point to proper file,
		// nothing to do here
		return nil
	}

	// file exists and is NOT a symbolic link

	// read dest/src files to determine if they have the same content
	destContents, err := ioutil.ReadFile(dest)
	if err != nil {
		return err
	}

	srcContents, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	// if files have the same content
	if strings.Compare(string(destContents), string(srcContents)) == 0 {
		logger.PrintDebug("FILE has same content as thing. Replacing file with symlink")
		if err := os.Remove(dest); err != nil {
			return err
		}

		if err := os.Symlink(src, dest); err != nil {
			return err
		}
		return nil
	}

	// replace with link
	logger.PrintInfo("FILE %s exists, is not a symlink and has different content. (remove, skip) ", src)
	var input string
	_, err = fmt.Scanln(&input)
	if err != nil {
		return err
	}

	switch input {
	case "remove":
		if err := os.Remove(dest); err != nil {
			return err
		}

		if err := os.Symlink(src, dest); err != nil {
			return err
		}

	case "skip":
		logger.PrintDebug("skipping '%s'\n", rel)
	}

	return nil
}

func resolveDirectory(src string, dest string, rel string) {

}

type File struct {
	File string   `toml:"file"`
	Tags []string `toml:"tags"`
	Type string   `toml:"type" default:"file"`
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
	Short: "Change User dotfiles",
	Long:  `Change User dotfiles`,
	Run: func(cmd *cobra.Command, args []string) {
		// get data
		storeDir := cmd.Flag("store-dir").Value.String()
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

			logger.PrintDebug("src %s\n", src)
			logger.PrintDebug("rel %s\n", rel)
			logger.PrintDebug("dest %s\n", dest)

			for _, file := range userDotsConfig.Files {
				if strings.HasSuffix(src, file.File) {
					fmt.Println()
					logger.PrintInfo("Match: %s\n", file.File)

					if info.IsDir() && file.Type == "directory" {
						resolveDirectory(src, dest, rel)
					} else if !info.IsDir() && file.Type == "file" {
						resolveFile(src, dest, rel)
					} else if info.IsDir() && file.Type == "file" {
						logger.PrintWarning("'%s' is specified as a file, but at '%s', it is actually a directory\n", file.File, src)
					} else {

					}
				}
			}

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
