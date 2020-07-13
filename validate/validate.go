package validate

import (
	"fmt"
	"os"
	"path"

	"github.com/eankeen/globe/internal/util"
	"github.com/spf13/cobra"
)

// Validate command line arguments and directory structure
func Validate(cmd *cobra.Command, args []string) {
	storeLocation := cmd.Flag("store-dir").Value.String()
	storeLocation = util.CheckFileStore(storeLocation)

	// checkCoreFiles(storeLocation)
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

	coreConfig := util.ReadCoreConfig(storeLocation)
	fmt.Print(coreConfig)
}
