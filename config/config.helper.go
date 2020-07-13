package config

import (
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/eankeen/globe/internal/util"
	"gopkg.in/yaml.v2"
)

// GlobeConfig is the configuration file used to manage Globe. This is your `globe.toml` file
type GlobeConfig struct {
	Globe struct {
		License string `toml:"license"`
	} `toml:"globe"`
	Init struct {
		Holds []string `toml:"holds"`
	} `toml:"init"`
	Sync struct {
		Holds []string `toml:"holds"`
	} `toml:"sync"`
}

// FileEntryRaw has data about a single file that is meant to be bootstrapped. It's raw because it comes stragith from the bootstrapFiles.yml file
type FileEntryRaw struct {
	Path string `yaml:"path"`
	For  string `yaml:"for"`
	Op   string `yaml:"op"`
}

// FileEntry is the same as FileEntryRaw, except it has been processed
type FileEntry struct {
	SrcPath  string `yaml:"srcPath"`
	DestPath string `yaml:"destPath"`
	RelPath  string `yaml:"relPath"`
	Op       string `yaml:"op"`
	For      string `yaml:"for"`
}

type FileListRaw struct {
	Files []FileEntryRaw `yaml:"files"`
}

type FileList struct {
	Files []FileEntry `yaml:"files"`
}

// ReadSyncConfig reads the local sync.yml configuration file
func ReadSyncConfig(storeDir string, storeLocation string) FileListRaw {
	yamlLocation := path.Join(storeDir, "sync.yml")

	var coreConfig FileListRaw
	{
		content, err := ioutil.ReadFile(yamlLocation)
		if err != nil {
			panic(err)
		}

		if err := yaml.Unmarshal(content, &coreConfig); err != nil {
			panic(err)
		}
	}

	return coreConfig
}

// ReadGlobeConfig reads the local globe.toml config file
func ReadGlobeConfig(projectDir string) GlobeConfig {
	configLocation := path.Join(projectDir, "globe.toml")

	var globeConfig GlobeConfig
	{
		content, err := ioutil.ReadFile(configLocation)
		if err != nil {
			panic(err)
		}

		if _, err = toml.Decode(string(content), &globeConfig); err != nil {
			panic(err)
		}
	}

	return globeConfig
}

func walkupFor(startLocation string, filename string) string {
	dirContents, err := ioutil.ReadDir(startLocation)
	if err != nil {
		log.Printf("Could not read directory %s", startLocation)
		panic(err)
	}

	util.PrintDebug("Searching for '%s' in %s\n", filename, startLocation)
	for _, file := range dirContents {
		// util.Debug("dir: '%s'\n", file.Name())
		if file.Name() == filename {
			util.PrintDebug("Found '%s' in '%s\n", filename, startLocation)
			return startLocation
		}
	}
	if startLocation == "/" {
		return ""
	}

	return walkupFor(path.Dir(startLocation), filename)
}

func GetProjectLocation() string {
	start, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return walkupFor(start, "globe.toml")
}

// GetDataLocation gets the full path to the configuration file configuration file
func GetDataLocation() string {
	projectDir := GetProjectLocation()
	return path.Join(projectDir, "globe.toml")
}
