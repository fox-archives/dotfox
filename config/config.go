package config

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/eankeen/globe/internal/util"
	logger "github.com/eankeen/go-logger"
)

// Project includes all details of the current Project
type Project struct {
	ProjectDir string
	StoreDir   string
	UserDir    string
	Config     Config
	Files      []FileEntry
}

// GetData gets the config for all data related to project
func GetData(storeDir string) Project {
	projectDir := GetProjectDir()

	var project Project
	project.StoreDir = storeDir

	logger.Debug("projectDir: %s\n", projectDir)
	project.ProjectDir = projectDir

	project.Config = ReadConfig(project.ProjectDir)

	homedir, err := os.UserHomeDir()
	util.P(err)

	project.UserDir = homedir

	// CONVERT FILE LISTS
	do := func(fileListRaw []FileEntryRaw) []FileEntry {
		var fileList []FileEntry

		for _, file := range fileListRaw {
			file := FileEntry{
				Op:       file.Op,
				For:      file.For,
				Tags:     file.Tags,
				Usage:    file.Usage,
				SrcPath:  path.Join(storeDir, file.Path),
				DestPath: path.Join(projectDir, file.Path),
				RelPath:  file.Path,
			}
			fileList = append(fileList, file)
		}

		return fileList
	}

	syncFilesRaw := ReadFileConfig(storeDir, projectDir)
	project.Files = do(syncFilesRaw.Files)
	// logger.Debug("syncFiles: %+v\n", project.Files)

	return project
}

// File represents an entry in the `user.dots.toml` file
type File struct {
	File       string   `toml:"file"`
	Tags       []string `toml:"tags"`
	Type       string   `toml:"type"`
	Heuristic1 bool
	Heuristic2 bool
	Heuristic3 bool
}

// UserDotsConfig represents the `user.dots.toml` file
type UserDotsConfig struct {
	Files []File `toml:"files"`
}

// GetUserToml gets User (~) data
func GetUserToml(storeDir string) UserDotsConfig {
	projectConfig := filepath.Join(storeDir, "user.dots.toml")

	raw, err := ioutil.ReadFile(projectConfig)
	util.P(err)

	var userDotsConfig UserDotsConfig
	err = toml.Unmarshal(raw, &userDotsConfig)
	util.P(err)

	return userDotsConfig
}
