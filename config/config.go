package config

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/eankeen/dotty/internal/t"
	"github.com/eankeen/dotty/internal/util"
)

// File represents a file entry in a `*.dots.toml` file
type File struct {
	File       string   `toml:"file"`
	Mode       int      `toml:"mode"`
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
	Files   []File   `toml:"files"`
	Ignores []Ignore `toml:"ignores"`
}

// UserDotsConfig represents the `user.dots.toml` file
type UserDotsConfig struct {
	Files   []File   `toml:"files"`
	Ignores []Ignore `toml:"ignores"`
}

// LocalDotsConfig represents the `local.dots.toml` file
type LocalDotsConfig struct {
	Files []File `toml:"files"`
}

// DottyCfg gets the `dotty.toml` file
func DottyCfg(dotfilesDir string) t.DottyConfig {
	file := filepath.Join(dotfilesDir, "dotty.toml")
	raw, err := ioutil.ReadFile(file)
	util.HandleFsError(err)

	var cfg t.DottyConfig
	err = toml.Unmarshal(raw, &cfg)
	util.HandleError(err)

	return cfg
}

// SystemCfg gets system (/) config
func SystemCfg(dotfilesDir string) SystemDotsConfig {
	file := filepath.Join(dotfilesDir, DottyCfg(dotfilesDir).ConfigDir, "system.dots.toml")
	raw, err := ioutil.ReadFile(file)
	util.HandleFsError(err)

	var cfg SystemDotsConfig
	err = toml.Unmarshal(raw, &cfg)
	util.HandleError(err)

	return cfg
}

// UserCfg gets user (~) config
func UserCfg(dotfilesDir string) UserDotsConfig {
	file := filepath.Join(dotfilesDir, DottyCfg(dotfilesDir).ConfigDir, "user.dots.toml")
	raw, err := ioutil.ReadFile(file)
	util.HandleFsError(err)

	var cfg UserDotsConfig
	err = toml.Unmarshal(raw, &cfg)
	util.HandleError(err)

	return cfg
}

// LocalCfg gets local (.) config
func LocalCfg(dotfilesDir string) LocalDotsConfig {
	file := filepath.Join(dotfilesDir, DottyCfg(dotfilesDir).ConfigDir, "local.dots.toml")

	raw, err := ioutil.ReadFile(file)
	util.HandleFsError(err)

	var cfg LocalDotsConfig
	err = toml.Unmarshal(raw, &cfg)
	util.HandleFsError(err)

	return cfg
}
