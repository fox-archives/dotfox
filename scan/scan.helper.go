package scan

import (
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/eankeen/globe/internal/util"
	"gopkg.in/yaml.v2"
)

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

// getProjectLocation returns an absolute path to the directory containing a `globe.toml` file
func getProjectLocation() string {
	start, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	projectLocation := walkupFor(start, "globe.toml")
	return projectLocation
}

// GetConfigLocation gets the full path to the configuration file configuration file
func GetConfigLocation() string {
	return path.Join(getProjectLocation(), "globe.toml")
}

func readGlobeConfig(projectLocation string) GlobeConfig {
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

func readBootstrapFilesRaw(projectLocation string) BootstrapFilesRaw {
	yamlLocation := path.Join(util.Dirname(), "bootstrapFiles.yml")

	var bootstrapFilesRaw BootstrapFilesRaw
	// TODO: bootstrapFiles = util.ReadYaml(yamlLocation) etc.
	{
		content, err := ioutil.ReadFile(yamlLocation)
		if err != nil {
			panic(err)
		}

		if err := yaml.Unmarshal(content, &bootstrapFilesRaw); err != nil {
			panic(err)
		}
	}

	return bootstrapFilesRaw
}

// Transform the BootstrapEntryRaw to BootstrapRaw
func createBootstrapFilesFromRaw(bootstrapFilesRaw BootstrapFilesRaw, projectLocation string) BootstrapFiles {
	dirname := util.Dirname()

	var bootstrapFiles BootstrapFiles
	for _, oldFile := range bootstrapFilesRaw.OldFiles {
		oldFile := BootstrapEntry{
			SrcPath:  path.Join(dirname, "files", oldFile.Path),
			DestPath: path.Join(projectLocation, oldFile.Path),
			RelPath:  oldFile.Path,
			For:      oldFile.For,
		}
		bootstrapFiles.OldFiles = append(bootstrapFiles.OldFiles, oldFile)
	}
	for _, newFile := range bootstrapFilesRaw.NewFiles {
		newFile := BootstrapEntry{
			SrcPath:  path.Join(dirname, "files", newFile.Path),
			DestPath: path.Join(projectLocation, newFile.Path),
			RelPath:  newFile.Path,
			For:      newFile.For,
		}
		bootstrapFiles.NewFiles = append(bootstrapFiles.NewFiles, newFile)
	}

	return bootstrapFiles
}
