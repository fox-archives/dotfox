package util

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
)

// FileExists stops the program if the file does not exist
func FileExists(name string) (bool, error) {
	_, err := os.Stat(name)

	if err != nil {
		if os.IsNotExist(err) {
			return false, err
		}
		return true, err
	}

	return true, nil
}

// // EnsureFolderExists stops the program if the file does not exist
// func EnsureFolderExists(name string) {
// 	stat, err := os.Stat(name)

// 	if err != nil {
// 		if os.IsNotExist(err) {
// 			return false, err
// 		}
// 		return false, err
// 	}

// 	return stat.IsDir(), nil
// }

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

// BootstrapFilesRaw lists the attributes of each file to bootstrap
type BootstrapFilesRaw struct {
	Files []BootstrapEntryRaw `yaml:"files"`
}

// BootstrapFiles is the same as BootstrapFilesRaw except is uses the processed versions
type BootstrapFiles struct {
	Files []BootstrapEntry `yaml:"files"`
}

// ReadCoreConfig reads the core config
func ReadCoreConfig(storeLocation string) BootstrapFilesRaw {
	yamlLocation := path.Join(storeLocation, "core/core.yml")

	var coreConfig BootstrapFilesRaw
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

// ReadBootstrapFilesRaw reads the core.yml file
func ReadBootstrapFilesRaw(storeLocation string) BootstrapFilesRaw {
	yamlLocation := path.Join(storeLocation, "core/core.yml")

	var coreConfig BootstrapFilesRaw
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
	// TODO: globeConfig = util.ReadToml(configLocation) etc.
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
