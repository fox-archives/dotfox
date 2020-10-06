package config

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/eankeen/dotty/internal/util"
)

// File represents a file entry in a `*.dots.toml` file
type File struct {
	File       string   `toml:"file"`
	Tags       []string `toml:"tags"`
	Heuristic1 bool
	Heuristic2 bool
	Heuristic3 bool
}

// Ignore represents an ignore entry in a `*.dots.toml` file
type Ignore struct {
	File string `toml:"file"`
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

// IgnoreMatches determines a particular ignore entry matches
// returned string can either be "folder" or "file". if bool
// is false, it can also be empty ("")
func IgnoreMatches(src string, file Ignore) (bool, string) {
	lastChar := file.File[len(file.File)-1:]

	// if src is a folder
	if lastChar == "/" {
		allButLastChar := file.File[:len(file.File)-1]
		return strings.Contains(src, allButLastChar), "folder"
	}

	// if src is a file
	return strings.Contains(src, file.File), "file"
}

// SystemDotsConfig represents the `system.dots.toml` file
type SystemDotsConfig struct {
	Files []File `toml:"files"`
}

// GetSystemToml gets system (/) config
func GetSystemToml(dotDir string) SystemDotsConfig {
	projectConfig := GetSystemTomlPath(dotDir)

	raw, err := ioutil.ReadFile(projectConfig)
	util.HandleFsError(err)

	var systemDotsConfig SystemDotsConfig
	err = toml.Unmarshal(raw, &systemDotsConfig)
	util.HandleError(err)

	return systemDotsConfig
}

// GetSystemTomlPath gets location of system (/) config (system.dots.toml)
func GetSystemTomlPath(dotDir string) string {
	location := filepath.Join(dotDir, "config", "system.dots.toml")

	return location
}

// UserDotsConfig represents the `user.dots.toml` file
type UserDotsConfig struct {
	Files   []File   `toml:"files"`
	Ignores []Ignore `toml:"ignores"`
}

// GetUserToml gets user (~) config
func GetUserToml(dotDir string) UserDotsConfig {
	projectConfig := GetUserTomlPath(dotDir)

	raw, err := ioutil.ReadFile(projectConfig)
	util.HandleFsError(err)

	var userDotsConfig UserDotsConfig
	err = toml.Unmarshal(raw, &userDotsConfig)
	util.HandleError(err)

	return userDotsConfig
}

// GetUserTomlPath gets location of user (~) config (user.dots.toml)
func GetUserTomlPath(dotDir string) string {
	location := filepath.Join(dotDir, "config", "user.dots.toml")

	return location
}

// LocalDotsConfig represents the `local.dots.toml` file
type LocalDotsConfig struct {
	Files []File `toml:"files"`
}

// GetLocalToml gets local (.) config
func GetLocalToml(dotDir string) LocalDotsConfig {
	projectConfig := GetLocalTomlPath(dotDir)

	raw, err := ioutil.ReadFile(projectConfig)
	util.HandleFsError(err)

	var localDotsConfig LocalDotsConfig
	err = toml.Unmarshal(raw, &localDotsConfig)
	util.HandleFsError(err)

	return localDotsConfig
}

// GetLocalTomlPath gets location of local (.) config (local.dots.toml)
func GetLocalTomlPath(dotDir string) string {
	location := filepath.Join(dotDir, "config", "local.dots.toml")

	return location
}
