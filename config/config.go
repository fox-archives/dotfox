package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/eankeen/dotty/internal/util"
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

// GetSystemToml gets system (/) config
func GetSystemToml(storeDir string) SystemDotsConfig {
	projectConfig := GetSystemTomlPath(storeDir)

	raw, err := ioutil.ReadFile(projectConfig)
	util.P(err)

	var systemDotsConfig SystemDotsConfig
	err = toml.Unmarshal(raw, &systemDotsConfig)
	util.P(err)

	return systemDotsConfig
}

// GetSystemTomlPath gets location of system (/) config (system.dots.toml)
func GetSystemTomlPath(storeDir string) string {
	location := filepath.Join(storeDir, "config", "system.dots.toml")

	return location
}

// UserDotsConfig represents the `user.dots.toml` file
type UserDotsConfig struct {
	Files []File `toml:"files"`
}

// GetUserToml gets user (~) config
func GetUserToml(storeDir string) UserDotsConfig {
	projectConfig := GetUserTomlPath(storeDir)

	raw, err := ioutil.ReadFile(projectConfig)
	util.P(err)

	var userDotsConfig UserDotsConfig
	err = toml.Unmarshal(raw, &userDotsConfig)
	util.P(err)

	return userDotsConfig
}

// GetUserTomlPath gets location of user (~) config (user.dots.toml)
func GetUserTomlPath(storeDir string) string {
	location := filepath.Join(storeDir, "config", "user.dots.toml")

	return location
}

// LocalDotsConfig represents the `local.dots.toml` file
type LocalDotsConfig struct {
	Files []File `toml:"files"`
}

// GetLocalToml gets local (.) config
func GetLocalToml(storeDir string) LocalDotsConfig {
	projectConfig := GetLocalTomlPath(storeDir)

	raw, err := ioutil.ReadFile(projectConfig)
	util.P(err)

	var localDotsConfig LocalDotsConfig
	err = toml.Unmarshal(raw, &localDotsConfig)
	util.P(err)

	return localDotsConfig
}

// GetLocalTomlPath gets location of local (.) config (local.dots.toml)
func GetLocalTomlPath(storeDir string) string {
	location := filepath.Join(storeDir, "config", "local.dots.toml")

	return location
}

// FileMatches determines a particular file matches
// returned string can either be "folder" or "file". if bool
// is false, it can also be empty ("")
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
