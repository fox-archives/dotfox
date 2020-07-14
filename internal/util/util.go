package util

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
)

// Dirname performs same function as `__dirname()` in Node, obtaining the parent folder of the file of the callee of this function
func Dirname() string {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		panic("could not recover information from call stack")
	}

	dir := path.Dir(filename)
	return dir
}

// GetChildFilesRecurse walks all the child files of a directory and returns them
func GetChildFilesRecurse(dir string) ([]string, error) {
	stat, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			PrintError("File or folder '%s' does not exist. Exiting.\n", dir)
			os.Exit(1)
		}
		panic(err)
	}
	if !stat.IsDir() {
		PrintError("The file '%s' is not a directory. Exiting.\n", dir)
		os.Exit(1)
	}

	files := []string{}
	err = filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})

	return files, err
}

// Contains tests to see if a particular string is in an array
func Contains(arr []string, str string) bool {
	for _, el := range arr {
		if el == str {
			return true
		}
	}
	return false
}
