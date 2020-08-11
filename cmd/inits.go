package cmd

import (
	"io"
	"os"
	"path"

	"github.com/eankeen/globe/internal/util"
	"github.com/spf13/cobra"
)

var initsCmd = &cobra.Command{
	Use:   "init",
	Short: "Init Globe's configuration files",
	Long:  `Initiates configuration files to be used by Globe`,
	Run: func(cmd *cobra.Command, args []string) {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		// COPY GLOBE.TOML
		{
			storeDir := cmd.Flag("store-dir").Value.String()
			srcConfig := path.Join(storeDir, "globe.toml")
			destConfig := path.Join(wd, "globe.toml")
			util.PrintDebug("Copying '%s' to '%s'\n", srcConfig, destConfig)

			sourceFile, err := os.Open(srcConfig)
			if err != nil {
				panic(err)
			}
			defer sourceFile.Close()

			// Create new file
			newFile, err := os.OpenFile(destConfig, os.O_CREATE|os.O_RDWR|os.O_EXCL, 0644)
			if err != nil {
				if os.IsExist(err) {
					util.PrintWarning("Config file 'globe.toml' file already exists. Not overwriting\n")
					goto createGlobeFolder
				}
				panic(err)
			}
			defer newFile.Close()

			_, err = io.Copy(newFile, sourceFile)
			if err != nil {
				panic(err)
			}
		}

		// CREATE .GLOBE FOLDER
	createGlobeFolder:
		{
			globeDotDir := path.Join(wd, ".globe")
			err = os.MkdirAll(globeDotDir, 0755)
			if err != nil {
				if os.IsExist(err) {
					util.PrintWarning("Folder `.globe` already exists. Not overwriting\n")
					return
				}
				util.PrintInfo("Error when creating `.globe` folder. Exiting.")
				panic(err)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(initsCmd)
}
