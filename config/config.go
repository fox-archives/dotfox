package config

import (
	"path"

	"github.com/eankeen/globe/internal/util"
)

// Project includes all details of the current Project
type Project struct {
	ProjectLocation string
	StoreDir        string
	GlobeConfig     GlobeConfig
	BootstrapFiles  BootstrapFiles
}

// GetConfig gets the config for all data related to project
func GetConfig(storeDir string) Project {
	var project Project
	project.StoreDir = storeDir

	projectLocation := getProjectLocation()
	util.PrintDebug("projectLocation: %s", projectLocation)
	project.ProjectLocation = projectLocation

	globeConfig := ReadGlobeConfig(projectLocation)
	util.PrintDebug("globeConfig: %+v\n", globeConfig)
	project.GlobeConfig = globeConfig

	bootstrapFilesRaw := ReadSyncConfig(storeDir, projectLocation)
	util.PrintDebug("readSyncConfigRaw: %+v\n", bootstrapFilesRaw)

	bootstrapFiles := createBootstrapFilesFromRaw(storeDir, bootstrapFilesRaw, projectLocation)
	util.PrintDebug("bootstrapFiles: %+v\n", bootstrapFiles)
	project.BootstrapFiles = bootstrapFiles

	return project
}

// Transform the BootstrapEntryRaw to BootstrapRaw
func createBootstrapFilesFromRaw(storeDir string, bootstrapFilesRaw SyncConfigRaw, projectLocation string) BootstrapFiles {

	var bootstrapFiles BootstrapFiles
	for _, file := range bootstrapFilesRaw.Files {
		file := BootstrapEntry{
			SrcPath:  path.Join(storeDir, "sync", file.Path),
			DestPath: path.Join(projectLocation, file.Path),
			RelPath:  file.Path,
			Op:       file.Op,
			For:      file.For,
		}
		bootstrapFiles.Files = append(bootstrapFiles.Files, file)
	}

	return bootstrapFiles
}
