package cmd

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path"

	"github.com/eankeen/dotty/config"
	"github.com/eankeen/dotty/fs"
	logger "github.com/eankeen/go-logger"
	"github.com/spf13/cobra"
)

var localApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Synchronize local dotfiles",
	Long:  "Synchronize local dotfiles. You will be prompted on conflicts",
	Run: func(cmd *cobra.Command, args []string) {
		writeGlobeState()
	},
}

func init() {
	localCmd.AddCommand(localApplyCmd)
}

func panicIfFileDoesNotExit(file string) {
	doesExist, err := fs.FilePossiblyExists(file)
	if err != nil {
		logger.Error("There was an error determining if there is a '%s' file or folder in the project directory\n", file)
		log.Panicln(err)
	}
	if !doesExist {
		logger.Error("The file '%s' could not be found. Did you forget to init?\n", file)
		log.Panicln("panicing due to unexpected error")
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

	// CREATE STRUCT, CREATE JSON TEXT, AND WRITE TO DISK
	var globeState = &GlobeState{
		OwnerName:               "Edwin Kofler",
		OwnerEmail:              "24364012+eankeen@users.noreply.github.com",
		OwnerWebsite:            "https://edwinkofler.com",
		Vcs:                     "git",
		VcsRemoteUsername:       "eankeen",
		VcsRemoteRepositoryName: path.Base(projectDir),
	}

	jsonText, err := json.MarshalIndent(globeState, "", "\t")
	if err != nil {
		logger.Error("There was a problem marshalling\n")
		log.Panicln(err)
	}

	err = ioutil.WriteFile(globeStateFile, jsonText, 0644)
	if err != nil {
		logger.Error("Error writing the 'globe.state.json' file\n")
		log.Panicln(err)
	}
}
