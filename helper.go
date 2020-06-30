package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/sergi/go-diff/diffmatchpatch"
)

type File struct {
	Path string `json:"path"`
	For  string `json:"for"`
}

type Files struct {
	Files    []File `json: files`
	OldFiles []File `json: oldFiles`
}

func GetAbsolutePaths(relativePath string) struct {
	SrcPath  string
	DestPath string
} {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	destPath := path.Join(dir, relativePath)

	srcPath := path.Join(_dirname(), "files", relativePath)

	return struct {
		SrcPath  string
		DestPath string
	}{
		SrcPath:  srcPath,
		DestPath: destPath,
	}
}

func ShouldRemoveExistingFile(path string, relativePath string, destContents []byte, srcContents []byte) bool {
	printInfo("File '%s' is outdated. Replace it? (y/d/n): ", relativePath)
	r := bufio.NewReader(os.Stdin)
	c, err := r.ReadByte()
	if err != nil {
		panic(err)
	}

	if c == byte('Y') || c == byte('y') {
		printInfo("chose: yes\n")
		return true
	} else if c == byte('N') || c == byte('n') {
		printInfo("chose: no\n")
		return false
	} else if c == byte('D') || c == byte('d') {
		printInfo("chose: diff\n")
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(string(destContents), string(srcContents), true)
		fmt.Println(dmp.DiffPrettyText(diffs))
		return ShouldRemoveExistingFile(path, relativePath, destContents, srcContents)
	} else {
		return ShouldRemoveExistingFile(path, relativePath, destContents, srcContents)
	}
}

func CopyFile(file File) {
	abs := GetAbsolutePaths(file.Path)
	srcFile := abs.SrcPath
	destFile := abs.DestPath

	// ensure parent directory exists
	os.MkdirAll(path.Dir(destFile), 0755)

	srcContents, err := ioutil.ReadFile(srcFile)
	if err != nil {
		panic(err)
	}

	// prompt to remove preexisting file if it exists
	stat, err := os.Stat(destFile)
	if stat != nil {
		// if the file buffers are the same, return no need to copy
		destContents, err := ioutil.ReadFile(destFile)
		if err != nil {
			panic(err)
		}

		if bytes.Compare(srcContents, destContents) == 0 {
			printInfo("Skipping unchanged '" + file.Path + "' file\n")
			return
		}

		// file exists, we ask if we should remove file
		shouldRemove := ShouldRemoveExistingFile(destFile, file.Path, destContents, srcContents)
		if shouldRemove == false {
			return
		}
	}

	err = ioutil.WriteFile(destFile, srcContents, 0644)
	if err != nil {
		log.Fatal(err)
	}

	printInfo("Copying %s to %s\n", srcFile, destFile)
}

func RemoveFile(file File) {
	abs := GetAbsolutePaths(file.Path)
	destFile := abs.DestPath

	err := os.Remove(destFile)
	if err != nil {
		fmt.Printf("Error when trying to remove %s. Skipping file", destFile)
		log.Println(err)
		return
	}

	// fileExists, err := FileExists(destFile)
	// if err != nil {
	// 	fmt.Printf("Error when checking if %s exists.", file.Path)
	// 	fmt.Println(err)
	// 	return
	// }

	// if fileExists {

	// }
}
