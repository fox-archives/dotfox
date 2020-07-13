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
	SyncFiles       FileList
	InitFiles       FileList
}

// GetData gets the config for all data related to project
func GetData(projectDir, storeDir string) Project {
	var project Project
	project.StoreDir = storeDir

	util.PrintDebug("projectDir: %s", projectDir)
	project.ProjectLocation = projectDir

	project.GlobeConfig = ReadGlobeConfig(projectDir)
	util.PrintDebug("globeConfig: %+v\n", project.GlobeConfig)

	bootstrapFilesRaw := ReadSyncConfig(storeDir, projectDir)
	util.PrintDebug("readSyncConfigRaw: %+v\n", bootstrapFilesRaw)

	project.SyncFiles = CreateSyncFilesFromRaw(storeDir, bootstrapFilesRaw, projectDir)
	util.PrintDebug("bootstrapFiles: %+v\n", project.SyncFiles)

	return project
}

// CreateSyncFilesFromRaw the FileEntryRaw to BootstrapRaw
func CreateSyncFilesFromRaw(storeDir string, bootstrapFilesRaw FileListRaw, projectDir string) FileList {

	var bootstrapFiles FileList
	for _, file := range bootstrapFilesRaw.Files {
		file := FileEntry{
			SrcPath:  path.Join(storeDir, "sync", file.Path),
			DestPath: path.Join(projectDir, file.Path),
			RelPath:  file.Path,
			Op:       file.Op,
			For:      file.For,
		}
		bootstrapFiles.Files = append(bootstrapFiles.Files, file)
	}

	return bootstrapFiles
}
