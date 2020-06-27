package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"

	"github.com/mattn/go-isatty"
)

type File struct {
	name string
}

func printInfo(text string, args ...interface{}) {
	_, isNoColorEnabled := os.LookupEnv("NO_COLOR")
	if (os.Getenv("TERM") != "dumb") && !isNoColorEnabled &&
		isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		fmt.Print("\033[32m")
		fmt.Printf(text, args...)
		fmt.Print("\033[m")
	} else {
		fmt.Printf(text, args...)
	}
}

func GetFullPaths(relativePath string) struct {
	SrcPath  string
	DestPath string
} {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	destPath := path.Join(dir, relativePath)

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("no caller information")
	}
	srcPath := path.Join(path.Dir(currentFile), "files", relativePath)

	return struct {
		SrcPath  string
		DestPath string
	}{
		SrcPath:  srcPath,
		DestPath: destPath,
	}
}

func ShouldRemoveExistingFile(path string, relativePath string) bool {
	printInfo("file '%s' is either old or not needed. remove? (y/n): ", relativePath)
	r := bufio.NewReader(os.Stdin)
	c, err := r.ReadByte()
	if err != nil {
		panic(err)
	}

	if c == byte('Y') || c == byte('y') {
		printInfo("chosen: yes\n")
		return true
	} else if c == byte('N') || c == byte('n') {
		printInfo("chosen: no\n")
		return false
	} else {
		return ShouldRemoveExistingFile(path, relativePath)
	}
}

func CopyFile(srcFile string, destFile string, relativePath string) {
	// ensure directory exists
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
			return
		}

		// file exists, we ask if we should remove file
		shouldRemove := ShouldRemoveExistingFile(destFile, relativePath)
		if shouldRemove == false {
			return
		}
	}

	err = ioutil.WriteFile(destFile, srcContents, 0644)
	if err != nil {
		log.Fatal(err)
	}

	printInfo("copying %s to %s\n", srcFile, destFile)
}

func GetFilesThatWereReplaced() []File {
	return []File{
		{
			name: "README.md",
		},
	}
}

func GetFilesToCopy() []File {
	return []File{
		{
			name: ".editorconfig",
		},
		{
			name: ".gitattributes",
		},
		{
			name: "license",
		},
		{
			name: ".github/CODE_OF_CONDUCT.md",
		},
		{
			name: ".github/COMMIT_CONVENTIONS.md",
		},
		{
			name: ".github/CONTRIBUTING.md",
		},
		{
			name: ".github/PULL_REQUEST_TEMPLATE.md",
		},
		{
			name: ".github/ISSUE_TEMPLATE/BUG_REPORT.md",
		},
		{
			name: ".github/ISSUE_TEMPLATE/config.yml",
		},
		{
			name: ".github/ISSUE_TEMPLATE/FEATURE_REQUEST.md",
		},
		{
			name: ".github/ISSUE_TEMPLATE/QUESTION.md",
		},
	}
}
