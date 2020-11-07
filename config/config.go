package config

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/eankeen/dotty/internal/t"
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
	file := GetCfgFile("system", dotfilesDir)
	raw, err := ioutil.ReadFile(file)
	util.HandleFsError(err)

	var cfg SystemDotsConfig
	err = toml.Unmarshal(raw, &cfg)
	util.HandleError(err)

	return cfg
}

// UserCfg gets user (~) config
func UserCfg(dotfilesDir string) UserDotsConfig {
	file := GetCfgFile("user", dotfilesDir)
	raw, err := ioutil.ReadFile(file)
	util.HandleFsError(err)

	var cfg UserDotsConfig
	err = toml.Unmarshal(raw, &cfg)
	util.HandleError(err)

	return cfg
}

// LocalCfg gets local (.) config
func LocalCfg(dotfilesDir string) LocalDotsConfig {
	file := GetCfgFile("local", dotfilesDir)
	raw, err := ioutil.ReadFile(file)
	util.HandleFsError(err)

	var cfg LocalDotsConfig
	err = toml.Unmarshal(raw, &cfg)
	util.HandleFsError(err)

	return cfg
}

// GetCfgFile gets the location of a particular dotfile (/, ~, or .) (system, user, or local)
func GetCfgFile(typ string, dotfilesDir string) string {
	configDir := DottyCfg(dotfilesDir).ConfigDir

	switch typ {
	case "system":
		location := filepath.Join(dotfilesDir, configDir, "system.dots.toml")
		return location

	case "user":
		location := filepath.Join(dotfilesDir, configDir, "user.dots.toml")
		return location

	case "local":
		location := filepath.Join(dotfilesDir, configDir, "local.dots.toml")
		return location

	default:
		log.Panicf("'%s' is not a valid type for GetConfigPath", typ)
		break
	}

	log.Panicf("'%s' is not a valid type for GetConfigPath", typ)
	return ""
}
