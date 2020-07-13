package sync

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/eankeen/globe/internal/util"
	"github.com/eankeen/globe/scan"
	"github.com/gobwas/glob"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func getAbsolutePaths(absolutePath string, relativePath string) struct {
	SrcPath  string
	DestPath string
} {
	dir := absolutePath
	destPath := path.Join(dir, relativePath)

	srcPath := path.Join(util.Dirname(), "files", relativePath)

	return struct {
		SrcPath  string
		DestPath string
	}{
		SrcPath:  srcPath,
		DestPath: destPath,
	}
}

func shouldRemoveExistingFile(path string, relativePath string, destContents []byte, srcContents []byte) bool {
	util.PrintInfo("FileEntry '%s' is outdated. Replace it? (y/d/n): ", relativePath)
	r := bufio.NewReader(os.Stdin)
	c, err := r.ReadByte()
	if err != nil {
		panic(err)
	}

	if c == byte('Y') || c == byte('y') {
		util.PrintInfo("chose: yes\n")
		return true
	} else if c == byte('N') || c == byte('n') {
		util.PrintInfo("chose: no\n")
		return false
	} else if c == byte('D') || c == byte('d') {
		util.PrintInfo("chose: diff\n")
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(string(destContents), string(srcContents), true)
		fmt.Println(dmp.DiffPrettyText(diffs))
		return shouldRemoveExistingFile(path, relativePath, destContents, srcContents)
	} else {
		return shouldRemoveExistingFile(path, relativePath, destContents, srcContents)
	}
}

func projectFilesContain(project scan.Project, glob glob.Glob) bool {
	dir := path.Join(util.Dirname(), "../scan/files")
	files, err := util.GetAllChildFolders(dir)
	if err != nil {
		panic(err)
	}

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

func isFileRelevant(project scan.Project, file util.BootstrapEntry) bool {
	projectContainsGoFiles := func() bool {
		if projectFilesContain(project, glob.MustCompile("*.go")) {
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

// CopyFile copies a file
func copyFile(project scan.Project, file util.BootstrapEntry) {
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
		fmt.Printf("Error trying to test if '%s' exists. Skipping file", destFile)
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
func removeFile(project scan.Project, file util.BootstrapEntry) {
	destFile := file.DestPath

	err := os.Remove(destFile)
	if err != nil {
		// fmt.Printf("Error when trying to remove %s. Skipping file\n", destFile)
		// log.Println(err)
		return
	}
}
