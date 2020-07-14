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
		// get data
		storeDir := cmd.Flag("store-dir").Value.String()
		srcConfig := path.Join(storeDir, "globe.toml")

		projectDir, err := os.Getwd()
		destConfig := path.Join(projectDir, "globe.toml")
		if err != nil {
			panic(err)
		}

		// COPY FILE
		{
			sourceFile, err := os.Open(srcConfig)
			if err != nil {
				panic(err)
			}
			defer sourceFile.Close()

			// Create new file
			newFile, err := os.OpenFile(destConfig, os.O_CREATE|os.O_EXCL, 0644)
			if err != nil {
				if os.IsExist(err) {
					util.PrintError("Config file 'globe.toml' file already exists. Not overwriting\n")
					return
				}
				panic(err)
			}
			defer newFile.Close()

			_, err = io.Copy(newFile, sourceFile)
			if err != nil {
				panic(err)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(initsCmd)
}
