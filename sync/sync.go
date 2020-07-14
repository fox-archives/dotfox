package sync

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/eankeen/globe/config"
	"github.com/eankeen/globe/internal/util"
)

func ProcessFiles(project config.Project, files []config.FileEntry) {
	for _, file := range files {
		util.PrintInfo("Processing file %s\n", file.RelPath)

		if file.Op == "add" {
			copyFile(project, file)
			continue
		} else if file.Op == "remove" {
			removeFile(project, file)
			continue
		}

		util.PrintError("File '%s's operation could not be read. Exiting.\n", file.RelPath)
	}
}

// CopyFile copies a file
func copyFile(project config.Project, file config.FileEntry) {
	srcFile := file.SrcPath
	destFile := file.DestPath
	util.PrintDebug("srcFile: %s\n", srcFile)
	util.PrintDebug("destFile: %s\n", destFile)

	// ensure parent directory exists
	os.MkdirAll(path.Dir(destFile), 0755)

	srcContents, err := ioutil.ReadFile(srcFile)
	if err != nil {
		panic(err)
	}

	// validate to see if we should even be trying to copy the file
	// over. for example scripts/go.sh should only be copied when
	// there are .go files in the repository
	isFileRelevant := isFileRelevant(project, file)
	if !isFileRelevant {
		util.PrintInfo("Non-relevant file '%s' is being skipped\n", file.RelPath)
		return
	}

	// prompt to remove preexisting file if it exists
	destFileExists, err := util.FileExists(destFile)
	if err != nil {
		fmt.Printf("Error trying to test if '%s' exists. Skipping file\n", destFile)
		log.Println(err)
		return
	}

	util.PrintDebug("destFileExists: %v\n", destFileExists)
	if destFileExists {
		// if the file buffers are the same, return no need to copy
		destContents, err := ioutil.ReadFile(destFile)
		if err != nil {
			panic(err)
		}

		if bytes.Compare(srcContents, destContents) == 0 {
			util.PrintInfo("Skipping unchanged '" + file.RelPath + "' file\n")
			return
		}

		// file exists, we ask if we should remove file
		shouldRemove := shouldRemoveExistingFile(destFile, file.RelPath, destContents, srcContents)
		if shouldRemove == false {
			return
		}
	}

	err = ioutil.WriteFile(destFile, srcContents, 0644)
	if err != nil {
		log.Fatal(err)
	}

	util.PrintInfo("Copying %s to %s\n", srcFile, destFile)
}

// RemoveFile removes a file
func removeFile(project config.Project, file config.FileEntry) {
	destFile := file.DestPath

	err := os.Remove(destFile)
	if err != nil {
		// fmt.Printf("Error when trying to remove %s. Skipping file\n", destFile)
		// log.Println(err)
		return
	}
}
