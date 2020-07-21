package config

import (
	"path"

	"github.com/eankeen/globe/internal/util"
	"github.com/spf13/cobra"
)

// Project includes all details of the current Project
type Project struct {
	ProjectDir  string
	StoreDir    string
	GlobeConfig GlobeConfig
	SyncFiles   FileList
}

// GetData gets the config for all data related to project
func GetData(cmd *cobra.Command, storeDir string) Project {
	projectDir := GetProjectDir()

	var project Project
	project.StoreDir = storeDir

	util.PrintDebug("projectDir: %s\n", projectDir)
	project.ProjectDir = projectDir

	// if we're not initiating, we read the global config
	if cmd.Name() != "init" {
		project.GlobeConfig = ReadGlobeConfig(projectDir)
		util.PrintDebug("globeConfig: %+v\n", project.GlobeConfig)
	}

	// CONVERT FILE LISTS
	do := func(fileListRaw FileListRaw) FileList {
		var fileList FileList
		for _, file := range fileListRaw.Files {
			file := FileEntry{
				SrcPath:  path.Join(storeDir, cmd.Name(), file.Path),
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
