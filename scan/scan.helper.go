package scan

import (
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/eankeen/globe/internal/util"
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

// Transform the BootstrapEntryRaw to BootstrapRaw
func createBootstrapFilesFromRaw(bootstrapFilesRaw util.BootstrapFilesRaw, projectLocation string) util.BootstrapFiles {
	dirname := util.Dirname()

	var bootstrapFiles util.BootstrapFiles
	for _, file := range bootstrapFilesRaw.Files {
		file := util.BootstrapEntry{
			SrcPath:  path.Join(dirname, "files", file.Path),
			DestPath: path.Join(projectLocation, file.Path),
			RelPath:  file.Path,
			Op:       file.Op,
			For:      file.For,
		}
		bootstrapFiles.Files = append(bootstrapFiles.Files, file)
	}

	return bootstrapFiles
}
