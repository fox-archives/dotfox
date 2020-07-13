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
	Bootstrap struct {
		Holds []string `toml:"holds"`
	} `toml:"bootstrap"`
}

// BootstrapEntryRaw has data about a single file that is meant to be bootstrapped. It's raw because it comes stragith from the bootstrapFiles.yml file
type BootstrapEntryRaw struct {
	Path string `yaml:"path"`
	For  string `yaml:"for"`
	Op   string `yaml:"op"`
}

// BootstrapEntry is the same as BootstrapEntryRaw, except it has been processed
type BootstrapEntry struct {
	SrcPath  string `yaml:"srcPath"`
	DestPath string `yaml:"destPath"`
	RelPath  string `yaml:"relPath"`
	Op       string `yaml:"op"`
	For      string `yaml:"for"`
}

// SyncConfigRaw lists the attributes of each file to bootstrap
type SyncConfigRaw struct {
	Files []BootstrapEntryRaw `yaml:"files"`
}

// BootstrapFiles is the same as SyncConfigRaw except is uses the processed versions
type BootstrapFiles struct {
	Files []BootstrapEntry `yaml:"files"`
}

// ReadSyncConfig reads the local sync.yml configuration file
func ReadSyncConfig(storeLocation string) SyncConfigRaw {
	yamlLocation := path.Join(storeLocation, "sync/sync.yml")

	var coreConfig SyncConfigRaw
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
func ReadGlobeConfig(projectLocation string) GlobeConfig {
	configLocation := path.Join(projectLocation, "globe.toml")

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

func getProjectLocation() string {
	start, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return walkupFor(start, "globe.toml")
}

// GetConfigLocation gets the full path to the configuration file configuration file
func GetConfigLocation() string {
	projectLocation := getProjectLocation()
	return path.Join(projectLocation, "globe.toml")
}
