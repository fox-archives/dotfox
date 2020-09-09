package config

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/eankeen/globe/internal/util"
	"gopkg.in/yaml.v2"
)

// GlobeConfig is the configuration file used to manage Globe. This is your `globe.toml` file
type GlobeConfig struct {
	Project struct {
		License string   `toml:"license"`
		Tags    []string `toml:"tags"`
	} `toml:"project"`
}

// FileEntryRaw has data about a single file that is meant to be bootstrapped. It's raw because it comes straight from the bootstrapFiles.yml file
type FileEntryRaw struct {
	Path  string   `yaml:"path"`
	For   []string `yaml:"for"`
	Op    string   `yaml:"op"`
	Tags  []string `yaml:"tags"`
	Usage string   `yaml:"usage"`
}

// FileEntry is the same as FileEntryRaw, except it has been processed
type FileEntry struct {
	SrcPath  string   `yaml:"srcPath"`
	DestPath string   `yaml:"destPath"`
	RelPath  string   `yaml:"relPath"`
	Op       string   `yaml:"op"`
	For      []string `yaml:"for"`
	Tags     []string `yaml:"tags"`
	Usage    string   `yaml:"usage"`
}

// FileListRaw is a representation of files to transfer
type FileListRaw struct {
	Files []FileEntryRaw `yaml:"files"`
}

// FileList is a representation of transformed files to transfer
type FileList struct {
	Files []FileEntry `yaml:"files"`
}

// ReadSyncConfig reads the local sync.yml configuration file
func ReadSyncConfig(storeDir string, storeLocation string) FileListRaw {
	// todo
	yamlLocation := filepath.Join(storeDir, "../project.sync.yml")

	var coreConfig FileListRaw
	content, err := ioutil.ReadFile(yamlLocation)
	if err != nil {
		util.PrintError("Could not read sync.yml located at '%s'. Exiting\n", yamlLocation)
		panic(err)
	}

	if err := yaml.Unmarshal(content, &coreConfig); err != nil {
		util.PrintError("Could not parse sync.yml located at '%s' as valid yaml. Exiting\n", yamlLocation)
		panic(err)
	}

	return coreConfig
}

// ReadGlobeConfig reads the local globe.toml config file
func ReadGlobeConfig(projectDir string) GlobeConfig {
	configLocation := path.Join(projectDir, "globe.toml")

	var globeConfig GlobeConfig
	content, err := ioutil.ReadFile(configLocation)
	if err != nil {
		panic(err)
	}

	if _, err = toml.Decode(string(content), &globeConfig); err != nil {
		panic(err)
	}

	return globeConfig
}

// GetProjectDir gets the root location of the current project, by recursively walking up directory tree until a globe.toml file is found. It stop searching after it reaches the user's home directory
func GetProjectDir() string {
	start, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return walkupFor(start, "globe.toml")
}

func walkupFor(startLocation string, filename string) string {
	dirContents, err := ioutil.ReadDir(startLocation)
	if err != nil {
		log.Printf("Could not read directory %s", startLocation)
		panic(err)
	}

	util.PrintDebug("Searching for '%s' in %s\n", filename, startLocation)
	for _, file := range dirContents {
		util.PrintDebug("dir: '%s'\n", file.Name())

		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}

		if file.Name() == filename {
			util.PrintDebug("Found '%s' in '%s\n", filename, startLocation)
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
