package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os/exec"
	"path"

	"github.com/eankeen/globe/config"
	"github.com/eankeen/globe/fs"
	"github.com/eankeen/globe/internal/util"
	"github.com/eankeen/globe/validate"
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

		// valudate values
		validate.Validate(validate.ValidationValues{
			StoreDir: storeDir,
			Project:  project,
		})

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

// GlobeState is the per-project state stored in the `globe.state` file
type GlobeState struct {
	OwnerFullname    string `json:"ownerFullname"`
	RepositoryRemote string `json:"repositoryRemote"`
	// RepositoryOwner  string `json:"repositoryOwner"`
	// RepositoryName   string `json:"repositoryName"`
}

func writeGlobeState() {
	projectDir := config.GetProjectDir()

	globeDotFolder := path.Join(projectDir, ".globe")
	doesExist, err := fs.FilePossiblyExists(globeDotFolder)
	if err != nil {
		util.PrintError("There was an error determining if there is a `.globe` folder in the project directory\n")
		panic(err)
	}
	if !doesExist {
		util.PrintError("The golder '.globe' could not be found in the current directory. Did you forget to init?\n")
		panic("")
	}

	cmd := exec.Command("git", "remote", "get-url", "origin")
	out, err := cmd.CombinedOutput()
	if err != nil {
		util.PrintError("There was an error when trying to get repository owner\n")
		panic(err)
	}

	globeStateDir := path.Join(globeDotFolder, "globe.state")
	var globeState = &GlobeState{
		OwnerFullname:    "Edwin Kofler",
		RepositoryRemote: string(out),
		// RepositoryOwner: ,
		// RepositoryName:  repositoryName,
	}
	globeStateText, err := json.Marshal(globeState)
	if err != nil {
		util.PrintError("There was a problem marshalling\n")
		panic(err)
	}

	err = ioutil.WriteFile(globeStateDir, globeStateText, 0644)
	if err != nil {
		util.PrintError("Error writing the 'globe.state' file\n")
		panic(err)
	}
}
