package scan

import (
	"github.com/eankeen/globe/internal/util"
)

// Project includes all details of the current Project. All information should be found in one pass
type Project struct {
	ProjectLocation string
	GlobeConfig     util.GlobeConfig
	BootstrapFiles  util.BootstrapFiles
}

// Scan scans for all data related to project
func Scan() Project {
	var project Project
	projectLocation := getProjectLocation()
	util.PrintDebug("projectLocation: %s", projectLocation)
	project.ProjectLocation = projectLocation

	globeConfig := util.ReadGlobeConfig(projectLocation)
	util.PrintDebug("globeConfig: %+v\n", globeConfig)
	project.GlobeConfig = globeConfig

	bootstrapFilesRaw := util.ReadBootstrapFilesRaw(projectLocation)
	util.PrintDebug("readBootstrapFilesRaw: %+v\n", bootstrapFilesRaw)

	bootstrapFiles := createBootstrapFilesFromRaw(bootstrapFilesRaw, projectLocation)
	util.PrintDebug("bootstrapFiles: %+v\n", bootstrapFiles)
	project.BootstrapFiles = bootstrapFiles

	return project
}
