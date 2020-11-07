package fs

import (
	"io/ioutil"
	"os"
	"path/filepath"

	logger "github.com/eankeen/go-logger"
)

// RemoveFile removes a file. If there are any errors in doing so, it immediately panics
func RemoveFile(destFile string) {
	err := os.Remove(destFile)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}

		logger.Error("Error when trying to remove file '%s'. Exiting\n", destFile)
		panic(err)
	}
}

// FilePossiblyExists stops the program if the file possiblyExists
// If no error is returned, we can be certain that boolean has
// integrity. If there is an error, then the file _possibly_ exists
// and the boolean does _not_ have integrity
func FilePossiblyExists(fileName string) (bool, error) {
	_, err := os.Stat(fileName)

	if err != nil {
		if os.IsNotExist(err) {
			// return nil because is a known error
			// that the value of the boolean depends on
			return false, nil
		}
		return true, err
	}
	return true, nil
}

// MkdirThenSymlink creates a new symlink to a destination. it
// automatically creates the parent directory structure too
func MkdirThenSymlink(src string, dest string) error {
	err := os.MkdirAll(filepath.Dir(dest), 0755)
	if err != nil {
		return err
	}

	err = os.Symlink(src, dest)
	if err != nil {
		return err
	}

	return nil
}

// RemoveThenSymlink removes a symlink that points to a wrong
// location, replacing it with the right one
func RemoveThenSymlink(src string, dest string) error {
	err := os.Remove(dest)
	if err != nil {
		return err
	}

	err = os.Symlink(src, dest)
	if err != nil {
		return err
	}

	return nil
}

// WriteTemp takes a `[]byte` and writes that to a temporary file, returning the File
func WriteTemp(content []byte) (os.File, error) {
	tempFile, err := ioutil.TempFile(os.TempDir(), "dotty-")
	if err != nil {
		return os.File{}, err
	}

	_, err = tempFile.Write(content)
	if err != nil {
		return os.File{}, err
	}

	return *tempFile, nil
}
