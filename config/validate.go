package config

import (
	"os"

	"github.com/eankeen/globe/internal/util"
)

// ValidatedArgs returns the value of cli arguments
type ValidatedArgs struct {
	StoreDir string
}

type ValidationValues struct {
	StoreDir string
	Project  Project `json:"omitempty"`
}

// Validate command line arguments and directory structure
func Validate(values ValidationValues) {
	// if store location is blank, we return prematurely
	// from this function because want cobra to print out
	// that the dot-dir is not set
	if values.StoreDir == "" {
		return
	}

	if values.Project.ProjectDir == "" {
		util.PrintError("projectDir is blank\n")
		panic("projectDir is blank")
	}

	// storeDir
	checkStoreDir(values.StoreDir)
}

func checkStoreDir(storeLocation string) {
	stat, err := os.Stat(storeLocation)

	if err != nil {
		if os.IsNotExist(err) {
			util.PrintError("The storeDir '%s'  does not exist. Exiting\n", storeLocation)
			os.Exit(1)
		}
		if os.IsPermission(err) {
			util.PrintError("There were permission issues when trying to stat '%s'. Exiting\n", storeLocation)
			os.Exit(1)
		}
		util.PrintError("An unknown error occured\n")
		panic(err)
	}

	if !stat.IsDir() {
		util.PrintError("Folder '%s' is not a folder. Exiting\n", storeLocation)
		os.Exit(1)
	}

	if storeLocation == "" {
		util.PrintError("fileStoreLocation is empty. This is not supposed to happen. Exiting\n")
		os.Exit(1)
	}
}
