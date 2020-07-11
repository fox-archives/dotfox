package inspect

import "github.com/eankeen/globe/internal/util"

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
}

// BootstrapEntry is the same as BootstrapEntryRaw, except it has been processed
type BootstrapEntry struct {
	SrcPath  string `yaml:"srcPath"`
	DestPath string `yaml:"destPath"`
	RelPath  string `yaml:"relPath"`
	For      string `yaml:"for"`
}

// BootstrapFilesRaw lists the attributes of each file to bootstrap
type BootstrapFilesRaw struct {
	NewFiles []BootstrapEntryRaw `yaml:"newFiles"`
	OldFiles []BootstrapEntryRaw `yaml:"oldFiles"`
}

// BootstrapFiles is the same as BootstrapFilesRaw except is uses the processed versions
type BootstrapFiles struct {
	NewFiles []BootstrapEntry `yaml:"newFiles"`
	OldFiles []BootstrapEntry `yaml:"oldFiles"`
}

// Project includes all details of the current Project. All information should be found in one pass
type Project struct {
	ProjectLocation string
	GlobeConfig     GlobeConfig
	BootstrapFiles  BootstrapFiles
}

// Inspect scans for all data related to project
func Inspect() Project {
	var project Project
	projectLocation := getProjectLocation()
	util.PrintDebug("projectLocation: %s", projectLocation)
	project.ProjectLocation = projectLocation

	globeConfig := readGlobeConfig(projectLocation)
	util.PrintDebug("globeConfig: %+v\n", globeConfig)
	project.GlobeConfig = globeConfig

	bootstrapFilesRaw := readBootstrapFilesRaw(projectLocation)
	util.PrintDebug("readBootstrapFilesRaw: %+v\n", bootstrapFilesRaw)

	bootstrapFiles := createBootstrapFilesFromRaw(bootstrapFilesRaw, projectLocation)
	util.PrintDebug("bootstrapFiles: %+v\n", bootstrapFiles)
	project.BootstrapFiles = bootstrapFiles

	return project
}
