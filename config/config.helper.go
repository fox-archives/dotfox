package config

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/eankeen/dotty/internal/util"
	logger "github.com/eankeen/go-logger"
)

// Config is the configuration file used to manage Globe. This is your `globe.toml` file
type Config struct {
	Project struct {
		License string   `toml:"license"`
		Tags    []string `toml:"tags"`
	} `toml:"project"`
}

// GetProjectDir gets the root location of the current project, by recursively walking up directory tree until a globe.toml file is found. It stop searching after it reaches the user's home directory
func GetProjectDir() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return walkupFor(wd, "globe.toml")
}

func walkupFor(startLocation string, filename string) string {
	dirContents, err := ioutil.ReadDir(startLocation)
	if err != nil {
		log.Printf("Could not read directory %s", startLocation)
		panic(err)
	}

	logger.Debug("Searching for '%s' in %s\n", filename, startLocation)
	for _, file := range dirContents {
		// logger.Debug("dir: '%s'\n", file.Name())

		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}

		if file.Name() == filename {
			logger.Debug("Found '%s' in '%s\n", filename, startLocation)
			return startLocation
		} else if file.Name() == homeDir {
			return ""
		}
	}
	if startLocation == "/" {
		return ""
	}

	return walkupFor(path.Dir(startLocation), filename)
}

// Src gets the location of a file, accounting for default values, config file values, and command line arguments
// TODO take into account command line arguments
func Src(dotfilesDir string, dottyCfg DottyConfig, typ string) string {
	switch typ {
	case "system":
		if dottyCfg.SystemDirSrc == "" {
			return filepath.Join(dotfilesDir, "system")
		}
		return pathExpand(dotfilesDir, dottyCfg.SystemDirSrc)
	case "user":
		if dottyCfg.UserDirSrc == "" {
			return filepath.Join(dotfilesDir, "user")
		}
		return pathExpand(dotfilesDir, dottyCfg.UserDirSrc)
	case "local":
		if dottyCfg.LocalDirSrc == "" {
			return filepath.Join(dotfilesDir, "local")
		}
		return pathExpand(dotfilesDir, dottyCfg.LocalDirSrc)
	}

	panic("Src not valid")
}

// Dest gets the location of a file, accounting for default values, config file values, and command line arguments
// TODO take into account command line arguments
func Dest(dotfilesDir string, dottyCfg DottyConfig, typ string) string {
	switch typ {
	case "system":
		if dottyCfg.SystemDirDest == "" {
			return "/"
		}
		return pathExpand(dotfilesDir, dottyCfg.SystemDirDest)
	case "user":
		if dottyCfg.UserDirDest == "" {
			homeDir, err := os.UserHomeDir()
			util.HandleFsError(err)

			return homeDir
		}
		return pathExpand(dotfilesDir, dottyCfg.UserDirDest)
	case "local":
		wd, err := os.Getwd()
		util.HandleError(err)

		return wd
	}

	panic("Dest not valid")
}

// PathExpand converts '~`, and to absolute path
func pathExpand(dotfilesDir string, rawPath string) string {
	isAbsolute := func(path string) bool {
		if strings.HasPrefix(path, "/") {
			return true
		}
		return false
	}

	if strings.HasPrefix(rawPath, "~") {
		homeDir, err := os.UserHomeDir()
		util.HandleFsError(err)
		rawPath = strings.Replace(rawPath, "~", homeDir, 1)
	}

	if strings.Contains(rawPath, "$HOME") {
		homeDir, err := os.UserHomeDir()
		util.HandleFsError(err)
		rawPath = strings.ReplaceAll(rawPath, "$HOME", homeDir)
	}

	if strings.Contains(rawPath, "$XDG_CONFIG_HOME") {
		configHome := os.Getenv("XDG_CONFIG_HOME")
		rawPath = strings.ReplaceAll(rawPath, "$XDG_CONFIG_HOME", configHome)
	}

	if isAbsolute(rawPath) {
		return rawPath
	}

	// relative
	return filepath.Join(dotfilesDir, rawPath)
}
