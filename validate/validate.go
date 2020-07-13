package validate

import (
	"os"
	"path"

	"github.com/eankeen/globe/internal/util"
)

// ValidatedArgs returns the value of cli arguments
type ValidatedArgs struct {
	StoreDir string
}

type ValidationValues struct {
	StoreDir string `json:"omitempty"`
}

// Validate command line arguments and directory structure
func Validate(values ValidationValues) {
	// if store location is blank, we return prematurelly
	// from this function because want cobra to print out
	// that the store-dir is not set
	if values.StoreDir == "" {
		return
	}

	// storeDir
	values.StoreDir = checkFileStore(values.StoreDir)

	// checkCoreFiles(storeDir)
}

func checkFileStore(storeLocation string) string {
	stat, err := os.Stat(storeLocation)

	if err != nil {
		if os.IsNotExist(err) {
			util.PrintError("The fileStore '%s'  does not exist. Exiting\n", storeLocation)
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

	return storeLocation
}

func checkCoreFiles(storeLocation string) {
	coreDir := path.Join(storeLocation, "core")
	// return if core directory cannot be found
	stat, err := os.Stat(coreDir)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		util.PrintError("An error occured\n")
		panic(err)
	}
	if !stat.IsDir() {
		util.PrintError("Folder '%s' is not a folder. Exiting\n", storeLocation)
	}

	// coreConfig := config.ReadSyncConfig(storeDir, storeLocation)
	// fmt.Print(coreConfig)
}
