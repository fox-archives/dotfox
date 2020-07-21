package config

import (
	"path"

	"github.com/eankeen/globe/internal/util"
)

// Project includes all details of the current Project
type Project struct {
	ProjectDir  string
	StoreDir    string
	GlobeConfig GlobeConfig
	SyncFiles   FileList
}

// GetData gets the config for all data related to project
func GetData(storeDir string) Project {
	projectDir := GetProjectDir()

	var project Project
	project.StoreDir = storeDir

	util.PrintDebug("projectDir: %s\n", projectDir)
	project.ProjectDir = projectDir

	// CONVERT FILE LISTS
	do := func(fileListRaw FileListRaw) FileList {
		var fileList FileList
		for _, file := range fileListRaw.Files {
			file := FileEntry{
				SrcPath:  path.Join(storeDir, "sync", file.Path),
				DestPath: path.Join(projectDir, file.Path),
				RelPath:  file.Path,
				Op:       file.Op,
				For:      file.For,
			}
			fileList.Files = append(fileList.Files, file)
		}

		return fileList
	}

	syncFilesRaw := ReadSyncConfig(storeDir, projectDir)

	project.SyncFiles = do(syncFilesRaw)

	util.PrintDebug("syncFiles: %+v\n", project.SyncFiles)

	return project
}
