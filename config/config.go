package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

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

// File represents an entry in the `user.dots.toml` file
type File struct {
	File       string   `toml:"file"`
	Tags       []string `toml:"tags"`
	Type       string   `toml:"type"`
	Heuristic1 bool
	Heuristic2 bool
	Heuristic3 bool
}

// SystemDotsConfig represents the `system.dots.toml` file
type SystemDotsConfig struct {
	Files []File `toml:"files"`
}

// GetSystemToml gets System (/) config
func GetSystemToml(storeDir string) SystemDotsConfig {
	projectConfig := filepath.Join(storeDir, "config", "system.dots.toml")

	raw, err := ioutil.ReadFile(projectConfig)
	util.P(err)

	var systemDotsConfig SystemDotsConfig
	err = toml.Unmarshal(raw, &systemDotsConfig)
	util.P(err)

	return systemDotsConfig
}

// UserDotsConfig represents the `user.dots.toml` file
type UserDotsConfig struct {
	Files []File `toml:"files"`
}

// GetUserToml gets user (~) config
func GetUserToml(storeDir string) UserDotsConfig {
	projectConfig := filepath.Join(storeDir, "config", "user.dots.toml")

	raw, err := ioutil.ReadFile(projectConfig)
	util.P(err)

	var userDotsConfig UserDotsConfig
	err = toml.Unmarshal(raw, &userDotsConfig)
	util.P(err)

	return userDotsConfig
}

// LocalDotsConfig represents the `local.dots.toml` file
type LocalDotsConfig struct {
	Files []File `toml:"files"`
}

// GetLocalToml gets local (.) config
func GetLocalToml(storeDir string) LocalDotsConfig {
	projectConfig := filepath.Join(storeDir, "config", "local.dots.toml")

	raw, err := ioutil.ReadFile(projectConfig)
	util.P(err)

	var localDotsConfig LocalDotsConfig
	err = toml.Unmarshal(raw, &localDotsConfig)
	util.P(err)

	return localDotsConfig
}

// FileMatches determines a particular file matches
// returned string can either be "folder" or "file"
func FileMatches(src string, file File) (bool, string) {
	lastChar := file.File[len(file.File)-1:]

	// if src is a folder
	if lastChar == "/" {
		allButLastChar := file.File[:len(file.File)-1]
		return strings.HasSuffix(src, allButLastChar), "folder"
	}

	// if src is a file
	return strings.HasSuffix(src, file.File), "file"
}

// CreateNewSymlink creates a new symlink to a destination. it
// automatically creates the parent directory structure too
func CreateNewSymlink(src string, dest string) error {
	logger.Debug("OK: dest '%s' doesn't exist. Recreating\n", dest)

	err := os.MkdirAll(filepath.Dir(dest), 0755)
	if err != nil {
		return err
	}

	err = os.Symlink(src, dest)
	if err != nil {
		return err
	}

	return nil
}

// FixBrokenSymlink removes a symlink that points to a wrong
// location, replacing it with the right one
func FixBrokenSymlink(src string, dest string) error {
	logger.Debug("OK: Symlink points to invalid location. Removing and Recreating\n")

	err := os.Remove(dest)
	if err != nil {
		return err
	}

	err = os.Symlink(src, dest)
	if err != nil {
		return err
	}

	return nil
}
