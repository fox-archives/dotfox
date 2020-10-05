package cmd

import (
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/eankeen/dotty/internal/util"
	logger "github.com/eankeen/go-logger"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init Dotty's config files",
	Long:  "Initializes Dotty's configuration files, usually located at ~/.config/dotty",
	Run: func(cmd *cobra.Command, args []string) {
		wd, err := os.Getwd()
		util.P(err)

		// COPY GLOBE.TOML
		{
			storeDir := cmd.Flag("dot-dir").Value.String()
			srcConfig := path.Join(storeDir, "globe.toml")
			destConfig := path.Join(wd, "globe.toml")
			logger.Debug("storeDir: %s\n", storeDir)
			logger.Debug("Copying '%s' to '%s'\n", srcConfig, destConfig)

			sourceFile, err := os.Open(srcConfig)
			defer sourceFile.Close()
			util.P(err)

			// Create new file
			newFile, err := os.OpenFile(destConfig, os.O_CREATE|os.O_RDWR|os.O_EXCL, 0644)
			if err != nil {
				if os.IsExist(err) {
					logger.Warning("Config file 'globe.toml' file already exists. Not overwriting\n")
					goto createGlobeFolder
				}
				panic(err)
			}
			defer newFile.Close()

			_, err = io.Copy(newFile, sourceFile)
			util.P(err)
		}

		// CREATE .GLOBE FOLDER
	createGlobeFolder:
		{
			globeDotDir := path.Join(wd, ".globe")
			err = os.MkdirAll(globeDotDir, 0755)
			if err != nil {
				if os.IsExist(err) {
					logger.Warning("Folder `.globe` already exists. Not overwriting\n")
					goto createGlobeStateJsonFile
				}
				logger.Informational("Error when creating `.globe` folder. Exiting.")
				panic(err)
			}
		}

		// CREATE GLOBE.STATE.JSON FILE
	createGlobeStateJsonFile:
		{
			globeStateJSONFile := path.Join(wd, ".globe", "globe.state.json")
			if ioutil.WriteFile(globeStateJSONFile, []byte("{}\n"), 0644); err != nil {
				if os.IsExist(err) {
					logger.Warning(("File .globe/globe.statea.json already exists. Not overwriting\n"))
					return
				}
				logger.Error("Could not create .globe/globe.state.json folder")
				panic(err)
			}
		}

	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}
