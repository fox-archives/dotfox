package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/mattn/go-isatty"
	"github.com/sergi/go-diff/diffmatchpatch"
)

type File struct {
	name string
}

func _dirname() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatalln("could not recover information from call stack")
	}

	dir := path.Dir(filename)
	return dir
}

func isColorEnabled() bool {
	_, isNoColorEnabled := os.LookupEnv("NO_COLOR")
	if (os.Getenv("TERM") != "dumb") && !isNoColorEnabled &&
		isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		return true
	}
	return false
}

func debug(text string, args ...interface{}) {
	isDebug := func() bool {
		_, ok := os.LookupEnv("DEBUG")
		if !ok {
			return false
		}
		return true
	}

	if isDebug() {
		if isColorEnabled() {
			fmt.Print("\033[32m")
			fmt.Printf(text, args...)
			fmt.Print("\033[m")
			return
		}

		fmt.Printf(text, args...)
	}
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

func GetFilesToCopyOver() []string {
	absoluteFiles, err := filepath.Glob(_dirname() + "/files/*")
	t, err := filepath.Glob(_dirname() + "/files/**/**")
	absoluteFiles = append(absoluteFiles, t...)

	if err != nil {
		log.Fatalln(err)
	}

	relativeFiles := []string{}
	for _, absoluteFile := range absoluteFiles {
		// skip directories
		{
			stat, err := os.Stat(absoluteFile)
			if err != nil {
				log.Fatalln(err)
			}
			if stat.IsDir() {
				debug("Skipping file: %s", absoluteFile)
				continue
			}
		}

		relativeFiles = append(relativeFiles, absoluteFile[len(_dirname()+"/files")+1:])
	}

	return relativeFiles
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

func CopyFile(srcFile string, destFile string, relativePath string) {
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
			printInfo("Skipping unchanged '" + relativePath + "' file\n")
			return
		}

		// file exists, we ask if we should remove file
		shouldRemove := ShouldRemoveExistingFile(destFile, relativePath, destContents, srcContents)
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

func GetFilesThatWereReplaced() []File {
	return []File{
		{
			name: "README.md",
		},
	}
}
