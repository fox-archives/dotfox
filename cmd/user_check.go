package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/eankeen/globe/config"
	"github.com/eankeen/globe/internal/util"
	"github.com/eankeen/go-logger"
	"github.com/spf13/cobra"
)

var userCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for inconsistencies / missing files",
	Run: func(cmd *cobra.Command, args []string) {
		storeDir := cmd.Flag("dot-dir").Value.String()
		userConfig := config.GetUserToml(storeDir)
		dotfileDir := filepath.Join(storeDir, "user")

		walkables := []string{}
		walkablesFound := []string{}
		err := filepath.Walk(dotfileDir, func(path string, info os.FileInfo, err error) error {
			// prevent errors in slice
			if path == dotfileDir {
				return nil
			}

			// TODO
			// homedir, err := os.UserHomeDir()
			// util.P(err)

			src := path
			// rel := path[len(dotfileDir)+1:]
			// dest := filepath.Join(homedir, rel)

			// for 3
			// TEMPORARILY ignore directories
			if !info.IsDir() {
				walkables = append(walkables, src)
			}

			// check to see if a file or folder is not used
			// heuristics:
			// 1. _toml file entry_ present, BUT not present in _dotfile repo_
			// 1. _toml folder entry_ present, BUT not present in _dotfile repo_
			// 2. _toml folder entry_ present, present in _dotfile repo_, BUT does not have any content (applying should fix this)
			for i, file := range userConfig.Files {
				if file.Type == "" {
					file.Type = "file"
				}

				// 1
				if config.FileMatches(src, file) {
					if info.IsDir() && file.Type == "folder" {
						userConfig.Files[i].Heuristic1 = true
					} else if !info.IsDir() && file.Type == "file" {
						userConfig.Files[i].Heuristic1 = true
					}
				}

				// 2
				if config.FileMatches(src, file) {
					if info.IsDir() && file.Type == "folder" {
						dirs, err := ioutil.ReadDir(src)
						util.P(err)

						if len(dirs) == 0 {
							userConfig.Files[i].Heuristic2 = true
						}
					}
				}

				// 3
				// TODO: recursively check parent directories, and check to make sure the
				if !strings.Contains(src, "fish/functions") || !strings.Contains(src, "oh-my-zsh") || !strings.Contains(src, "bash-it") {

					// file isn't covered in a parent _folder_ symlink
					if !config.FileMatches(src, file) {
						// TEMPORARILY ignore directories
						if !info.IsDir() {
							if !ParentFolderMatches(dotfileDir, info, src, userConfig.Files, file) {
								walkablesFound = append(walkablesFound, src)

							}

						}
					}

				}
			}

			return nil
		})
		util.P(err)

		// now display
		fmt.Println("Fails if you have an entry in your .toml file, but no file/folder was matched in your dotfile repo")
		for _, file := range userConfig.Files {
			if file.Heuristic1 == false {
				logger.Informational("Failed Heuristic1: %s\n", file.File)
			}
		}

		fmt.Println()
		fmt.Println("Fails if you have am empty folder in your dotfile repo")
		for _, file := range userConfig.Files {
			if file.Heuristic2 == true {
				logger.Informational("Failed Heuristic2: %s\n", file.File)
			}
		}

		fmt.Println()
		fmt.Println("Fails if you have a file in your dotfile repo, but it was not matched by anything in your toml file")
		// it is important that we loop over each walkable (representing a file in dotfile repo)
		for _, file := range walkables {
			if !util.Contains(walkablesFound, file) {
				logger.Informational("Failed Heuristic3: %s\n", file)
			}
		}
	},
}

func init() {
	userCmd.AddCommand(userCheckCmd)
}

// ParentFolderMatches to see if any parent folder of a file matches up until dotfileDir
func ParentFolderMatches(dotfileDir string, info os.FileInfo, src string, files []config.File, file config.File) bool {
	os := src

	if info.IsDir() {
		panic("not supposed to be directory")
	}

	for true {
		src = filepath.Dir(src)

		fmt.Println(dotfileDir, src)
		if dotfileDir == src || src == "/" {
			return false
		}

		for _, file := range files {
			if config.FileMatches(src, file) {
				fmt.Printf("CONTAINED: %s\n", os)
				return true
			}
		}
	}
	return false
}
