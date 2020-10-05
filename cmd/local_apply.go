package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/eankeen/globe/config"
	"github.com/eankeen/globe/fs"
	"github.com/eankeen/globe/internal/util"
	logger "github.com/eankeen/go-logger"
	"github.com/gobwas/glob"
	"github.com/spf13/cobra"
)

var localApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply config intelligently",
	Run: func(cmd *cobra.Command, args []string) {
		// write globe.state
		writeGlobeState()

		// get data
		storeDir := cmd.Flag("dot-dir").Value.String()
		projectDir := config.GetProjectDir()

		var project config.Project
		project.StoreDir = storeDir

		logger.Debug("projectDir: %s\n", projectDir)
		project.ProjectDir = projectDir

		project.Config = config.ReadConfig(project.ProjectDir)

		homedir, err := os.UserHomeDir()
		util.P(err)

		project.UserDir = homedir

		// CONVERT FILE LISTS
		do := func(fileListRaw []config.FileEntryRaw) []config.FileEntry {
			var fileList []config.FileEntry

			for _, file := range fileListRaw {
				file := config.FileEntry{
					Op:       file.Op,
					For:      file.For,
					Tags:     file.Tags,
					Usage:    file.Usage,
					SrcPath:  path.Join(storeDir, file.Path),
					DestPath: path.Join(projectDir, file.Path),
					RelPath:  file.Path,
				}
				fileList = append(fileList, file)
			}

			return fileList

		}

		syncFilesRaw := config.ReadFileConfig(storeDir, projectDir)
		project.Files = do(syncFilesRaw.Files)

		// process filesproject
		ProcessFiles(project, project.Files)
	},
}

func init() {
	localCmd.AddCommand(localApplyCmd)
}

// ProcessFiles processes each file according to their
// 'Op' operation specified by the user
func ProcessFiles(project config.Project, files []config.FileEntry) {
	for _, file := range files {
		logger.Informational("Processing file %s\n", file.RelPath)

		if file.Op == "add" {
			isFileRelevant := isFileRelevant(project, file)
			if !isFileRelevant {
				logger.Informational("Skipping irrelevant file '%s'\n", file.RelPath)
				continue
			}
			fs.CopyFile(file.SrcPath, file.DestPath, file.RelPath, project)
			continue
		} else if file.Op == "remove" {
			fs.RemoveFile(file.DestPath)
			continue
		}

		logger.Error("File '%s's operation could not be read. Exiting.\n", file.RelPath)
	}
}

func isFileRelevant(project config.Project, file config.FileEntry) bool {
	for _, tag := range file.Tags {
		if util.Contains(project.Config.Project.Tags, tag) {
			logger.Debug("tag: %s\n", tag)
			return true
		}
	}
	return false
}

func projectFilesContain(files []string, glob glob.Glob) bool {
	var doesContain bool
	for _, file := range files {
		if glob.Match(file) {
			doesContain = true
			break
		}
	}

	logger.Debug("Does project contain pattern %+v?: %t\n", glob, doesContain)
	return doesContain
}

func panicIfFileDoesNotExit(file string) {
	doesExist, err := fs.FilePossiblyExists(file)
	if err != nil {
		logger.Error("There was an error determining if there is a '%s' file or folder in the project directory\n", file)
		panic(err)
	}
	if !doesExist {
		logger.Error("The file '%s' could not be found. Did you forget to init?\n", file)
		panic("panicing due to unexpected error")
	}
}

// GlobeState is the per-project state stored in the `globe.state` file
type GlobeState struct {
	OwnerName               string `json:"ownerName"`
	OwnerEmail              string `json:"ownerEmail"`
	OwnerWebsite            string `json:"ownerWebsite"`
	Vcs                     string `json:"vcs"`
	VcsRemoteUsername       string `json:"vcsRemoteUsername"`
	VcsRemoteRepositoryName string `json:"vcsRemoteRepositoryName"`
}

func writeGlobeState() {
	projectDir := config.GetProjectDir()

	globeDotDir := path.Join(projectDir, ".globe")
	globeStateFile := path.Join(globeDotDir, "globe.state.json")

	panicIfFileDoesNotExit(globeDotDir)

	// OWNERNAME
	var ownerFullname string
	{
		ownerFullname = "Edwin Kofler"
	}

	// OWNEREMAIL
	var ownerEmail string
	{
		ownerEmail = "24364012+eankeen@users.noreply.github.com"
	}

	// // REPOSITORY REMOTE
	// var remoteRepo []byte
	// {
	// 	cmd := exec.Command("git", "remote", "get-url", "origin")
	// 	var err error
	// 	remoteRepo, err = cmd.CombinedOutput()
	// 	if err != nil {
	// 		logger.Error("There was an error when trying to get repository owner\n")
	// 		panic(err)
	// 	}
	// }

	// OWNERWEBSITE
	var ownerWebsite string
	{
		ownerWebsite = "https://edwinkofler.com"
	}

	// VCS
	var vcs string
	{
		vcs = "git"
	}

	// VCSREMOTEUSERNAME
	var vcsRemoteUsername string
	{
		vcsRemoteUsername = "eankeen"
	}

	// VCSREMOTEREPOSITORYNAME
	var vcsRemoteRepositoryName string
	{
		vcsRemoteRepositoryName = path.Base(projectDir)
	}

	// CREATE STRUCT, CREATE JSON TEXT, AND WRITE TO DISK
	var globeState = &GlobeState{
		OwnerName:               ownerFullname,
		OwnerEmail:              ownerEmail,
		OwnerWebsite:            ownerWebsite,
		Vcs:                     vcs,
		VcsRemoteUsername:       vcsRemoteUsername,
		VcsRemoteRepositoryName: vcsRemoteRepositoryName,
	}

	jsonText, err := json.MarshalIndent(globeState, "", "\t")
	if err != nil {
		logger.Error("There was a problem marshalling\n")
		panic(err)
	}

	err = ioutil.WriteFile(globeStateFile, jsonText, 0644)
	if err != nil {
		logger.Error("Error writing the 'globe.state.json' file\n")
		panic(err)
	}
}
