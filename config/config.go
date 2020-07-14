package config

import (
	"os"
	"path"

	"github.com/eankeen/globe/internal/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Project includes all details of the current Project
type Project struct {
	ProjectDir  string
	StoreDir    string
	GlobeConfig GlobeConfig
	SyncFiles   FileList
	InitFiles   FileList
}

// GetData gets the config for all data related to project
func GetData(cmd *cobra.Command, projectDir string, storeDir string) Project {
	var project Project
	project.StoreDir = storeDir

	// validate
	func() {
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				if cmd.Name() != "init" {
					util.PrintError("Please initiate a Globe project first\n")
					os.Exit(1)
				}
			}

			if os.IsNotExist(err) && cmd.Name() != "init" {
				util.PrintError("Please initiate a Globe project first\n")
				os.Exit(1)
			}

			if cmd.Name() != "init" {
				util.PrintError("An unknown error occured\n")
				panic(err)
			}
		}
	}()

	util.PrintDebug("projectDir: %s\n", projectDir)
	project.ProjectDir = projectDir

	// if we're not initiating, we read the global config
	if cmd.Name() != "init" {
		project.GlobeConfig = ReadGlobeConfig(projectDir)
		util.PrintError("globeConfig: %+v\n", project.GlobeConfig)
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
	initFilesRaw := ReadInitConfig(storeDir, projectDir)

	project.SyncFiles = do(syncFilesRaw)
	project.InitFiles = do(initFilesRaw)

	util.PrintDebug("syncFiles: %+v\n", project.SyncFiles)
	util.PrintDebug("syncFiles: %+v\n", project.InitFiles)

	return project
}
