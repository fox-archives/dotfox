package config

import (
	"os"

	logger "github.com/eankeen/go-logger"
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
		logger.Error("projectDir is blank\n")
		panic("projectDir is blank")
	}

	// storeDir
	checkStoreDir(values.StoreDir)
}

func checkStoreDir(storeLocation string) {
	stat, err := os.Stat(storeLocation)

	if err != nil {
		if os.IsNotExist(err) {
			logger.Error("The storeDir '%s'  does not exist. Exiting\n", storeLocation)
			os.Exit(1)
		}
		if os.IsPermission(err) {
			logger.Error("There were permission issues when trying to stat '%s'. Exiting\n", storeLocation)
			os.Exit(1)
		}
		logger.Error("An unknown error occured\n")
		panic(err)
	}

	if !stat.IsDir() {
		logger.Error("Folder '%s' is not a folder. Exiting\n", storeLocation)
		os.Exit(1)
	}

	if storeLocation == "" {
		logger.Error("fileStoreLocation is empty. This is not supposed to happen. Exiting\n")
		os.Exit(1)
	}
}
