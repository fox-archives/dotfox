package util

import (
	"os"
	"path"
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

// FileExists test if a file exists
// TODO: buggy
func FileExists(name string) (bool, error) {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
	}
	return true, nil
}
