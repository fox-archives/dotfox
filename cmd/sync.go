package cmd

import (
	"encoding/json"
	"io/ioutil"
	"path"

	"github.com/eankeen/globe/config"
	"github.com/eankeen/globe/fs"
	"github.com/eankeen/globe/internal/util"
	"github.com/gobwas/glob"
	"github.com/spf13/cobra"
)

var syncCommand = &cobra.Command{
	Use:   "sync",
	Short: "Sync Globe's configuration files",
	Long:  `Syncs configuration files`,
	Run: func(cmd *cobra.Command, args []string) {
		// write globe.state
		writeGlobeState()

		// get data
		storeDir := cmd.Flag("store-dir").Value.String()
		project := config.GetData(storeDir)

		// process files
		ProcessFiles(project, project.SyncFiles.Files)
	},
}

func init() {
	RootCmd.AddCommand(syncCommand)
}

// ProcessFiles processes each file according to their
// 'Op' operation specified by the user
func ProcessFiles(project config.Project, files []config.FileEntry) {
	for _, file := range files {
		util.PrintDebug("Processing file %s\n", file.RelPath)

		if file.Op == "add" {

			// validate to see if we should even be trying to copy the file
			// over. for example scripts/go.sh should only be copied when
			// there are .go files in the repository
			isFileRelevant := isFileRelevant(project.ProjectDir, file)
			if !isFileRelevant {
				util.PrintInfo("Skipping irrelevant file '%s'\n", file.RelPath)
				return
			}
			fs.CopyFile(file.SrcPath, file.DestPath, file.RelPath)
			continue
		} else if file.Op == "remove" {
			fs.RemoveFile(file.DestPath)
			continue
		}

		util.PrintError("File '%s's operation could not be read. Exiting.\n", file.RelPath)
	}
}

func isFileRelevant(projectDir string, file config.FileEntry) bool {
	projectContainsGoFiles := func() bool {
		files, err := util.GetChildFilesRecurse(projectDir)
		if err != nil {
			panic(err)
		}
		if projectFilesContain(files, glob.MustCompile("*.go")) {
			return true
		}
		return false

	}

	switch file.For {
	case "all":
		return true
	case "golang":
		if projectContainsGoFiles() {
			return true
		}
		return false
	}

	util.PrintDebug("FileEntry '%s' does not match case statement. Has value %s. Skipping\n", file.RelPath, file.For)
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

	util.PrintDebug("Does project contain pattern %+v?: %t\n", glob, doesContain)
	return doesContain
}

func panicIfFileDoesNotExit(file string) {
	doesExist, err := fs.FilePossiblyExists(file)
	if err != nil {
		util.PrintError("There was an error determining if there is a '%s' file in the project directory\n", file)
		panic(err)
	}
	if !doesExist {
		util.PrintError("The file '%s' could not be found. Did you forget to init?\n", file)
		panic("panicing due to unexpected error")
	}
}

// GlobeState is the per-project state stored in the `globe.state` file
type GlobeState struct {
	OwnerName               string `json:"ownerName"`
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

	// // REPOSITORY REMOTE
	// var remoteRepo []byte
	// {
	// 	cmd := exec.Command("git", "remote", "get-url", "origin")
	// 	var err error
	// 	remoteRepo, err = cmd.CombinedOutput()
	// 	if err != nil {
	// 		util.PrintError("There was an error when trying to get repository owner\n")
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
		OwnerWebsite:            ownerWebsite,
		Vcs:                     vcs,
		VcsRemoteUsername:       vcsRemoteUsername,
		VcsRemoteRepositoryName: vcsRemoteRepositoryName,
	}

	jsonText, err := json.MarshalIndent(globeState, "", "\t")
	if err != nil {
		util.PrintError("There was a problem marshalling\n")
		panic(err)
	}

	err = ioutil.WriteFile(globeStateFile, jsonText, 0644)
	if err != nil {
		util.PrintError("Error writing the 'globe.state' file\n")
		panic(err)
	}
}
