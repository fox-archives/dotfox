package config

import (
	"io/ioutil"
	"log"
	"os"
	"path"

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
