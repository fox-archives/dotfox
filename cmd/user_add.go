package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/eankeen/dotty/internal/util"
	"github.com/eankeen/go-logger"
	"github.com/spf13/cobra"
)

// assume if type is 'file' or 'folder'
func addFile(filePath string, fileType string) {
	if fileType == "file" {
		stat, err := os.Stat(filePath)
		// we get error if the file doesn't exist
		if err != nil {
			err := os.MkdirAll(filepath.Dir(filePath), 0755)
			util.HandleFsError(err)

			newFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_EXCL, 0644)
			util.HandleError(err)
			defer func() {
				err := newFile.Close()
				util.HandleFsError(err)
			}()

			logger.Informational("Creating file at '%s'\n", filePath)
			return
		}

		if stat.IsDir() {
			// TODO: automatically replace folder with file if it is empty
			logger.Error("'%s' exists and is a folder. Please remove folder before continuing\n", filePath)
			return
		}

		logger.Notice("'%s' already exists\n", filePath)
	} else if fileType == "folder" {
		stat, err := os.Stat(filePath)
		// we get error if the file doesn't exist
		if err != nil {
			logger.Informational("Creating folder at '%s'\n", filePath)
			err := os.MkdirAll(filePath, 0755)
			util.HandleFsError(err)
			return
		}

		if !stat.IsDir() {
			// TODO: automatically replace file with folder if it is empty
			logger.Error("'%s' exists and is a file. Please remove the file\n", filePath)
			return
		}

		logger.Notice("Folder '%s' already exists. Not doing anything\n", filePath)
	} else {
		logger.Error("fileType not either 'file' or 'folder'\n")
	}
}

var userAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a dotfile",
	Long:  "Adds a single dotfile",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			os.Stderr.WriteString("Error: Must add a path to add\n")
			os.Exit(1)
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		userDir := cmd.Flag("user-dir-dest").Value.String()
		dotfilesDir := cmd.Flag("dotfiles-dir").Value.String()

		fileType := cmd.Flag("type").Value.String()
		if !(fileType == "folder" || fileType == "file") {
			logger.Error("Type must either be 'file' or 'folder'")

		}

		newPath := args[0]

		// we place this before userDir because dotfilesDir may include
		// the contents of userDir, plus an extra path
		if strings.HasPrefix(newPath, dotfilesDir) {
			addFile(newPath, fileType)
			return
		}

		// test if has userDir prefix, replace it with dotfilesDir
		if strings.HasPrefix(newPath, userDir) {
			str := filepath.Join(dotfilesDir, "user", newPath[len(userDir):])
			addFile(str, fileType)
			return
		}

		if filepath.IsAbs(newPath) {
			logger.Error("Your path was not recognized. It must be either a relative path (relative to the dotfilesDir) or an absolute path that has a prefix of your dotfilesDir or userDir")
			return
		}

		// if path not absolute, assume it is relative to dotfilesDir
		newPath = filepath.Join(dotfilesDir, "user", newPath)
		addFile(newPath, fileType)
	},
}

func init() {
	userCmd.AddCommand(userAddCmd)
	pf := userAddCmd.PersistentFlags()
	pf.String("type", "", "Type of the newly added path")

	cobra.MarkFlagRequired(pf, "type")
}
