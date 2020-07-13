package util

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
)

// Dirname performs same function as `__dirname()` in Node
func Dirname() string {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		panic("could not recover information from call stack")
		// log.Fatalln("could not recover information from call stack")
	}

	dir := path.Dir(filename)
	return dir
}

// GetAllChildFolders walks all the child files of a directory and returns them
func GetAllChildFolders(dir string) ([]string, error) {
	stat, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			PrintError("File or folder '%s' does not exist. Exiting.", dir)
			os.Exit(1)
		}
		panic(err)
	}
	if !stat.IsDir() {
		PrintError("The file '%s' is not a directory. Exiting.", dir)
		os.Exit(1)
	}

	files := []string{}
	err = filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})

	return files, err
}

func CheckFileStore(storeLocation string) string {
	stat, err := os.Stat(storeLocation)

	if err != nil {
		if os.IsNotExist(err) {
			PrintError("The fileStore '%s'  does not exist. Exiting\n", storeLocation)
			os.Exit(1)
		}
		if os.IsPermission(err) {
			PrintError("There were permission issues when trying to stat '%s'. Exiting\n", storeLocation)
			os.Exit(1)
		}
		PrintError("An unknown error occured\n")
		panic(err)
	}

	if !stat.IsDir() {
		PrintError("Folder '%s' is not a folder. Exiting\n", storeLocation)
		os.Exit(1)
	}

	if storeLocation == "" {
		PrintError("fileStoreLocation is empty. This is not supposed to happen. Exiting\n")
		os.Exit(1)
	}

	return storeLocation
}
