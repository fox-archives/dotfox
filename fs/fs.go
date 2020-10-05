package fs

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"

	"github.com/eankeen/dotty/config"
	logger "github.com/eankeen/go-logger"
)

// CopyFile copies a file from a source to destination. If there are any errors,
// it prints the error to the screen and immediately panics
func CopyFile(srcFile string, destFile string, relFile string, templateVars config.Project) {
	logger.Debug("srcFile: %s\n", srcFile)
	logger.Debug("destFile: %s\n", destFile)

	// ensure parent directory exists
	err := os.MkdirAll(path.Dir(destFile), 0755)
	if err != nil {
		logger.Error("An error occurred when trying to recursively create a directory at '%s'. Exiting\n", destFile)
		panic(err)
	}

	srcContents, err := ioutil.ReadFile(srcFile)
	templatedSrcContents := templateFile(srcContents, templateVars, relFile)
	if err != nil {
		logger.Error("An error occurred when trying to read the file '%s'. Exiting\n", srcFile)
		panic(err)
	}

	// prompt to remove preexisting file if it exists
	destFilePossiblyExists, err := FilePossiblyExists(destFile)
	if err != nil {
		fmt.Printf("Could not determine if destination file '%s' exists. It could, but we received an error when trying to determine so. Exiting\n", destFile)
		panic(err)
	}
	// since we panic if there is an error, from now on we can
	// be certain that the boolean indicates if the file exist
	fileExists := destFilePossiblyExists

	logger.Debug("fileExists: %v\n", fileExists)

	// only continue if we are sure the destination file does not exist. of course, there can still be races, but we'll make sure to print errors
	if fileExists {
		// if the file buffers are the same, return no need to copy
		destContents, err := ioutil.ReadFile(destFile)
		if err != nil {
			logger.Error("An error occurred when trying to read the file '%s'. Exiting\n", destContents)
			panic(err)
		}

		// if the files are the same, don't copy and return
		if bytes.Compare(templatedSrcContents, destContents) == 0 {
			logger.Debug("Skipping unchanged '%s' file\n", relFile)
			return
		}

		// file exists and are different, we ask if we should remove file
		shouldRemove := shouldRemoveExistingFile(destFile, relFile, destContents, templatedSrcContents)
		if shouldRemove == false {
			return
		}
	}

	// if we got here, it means the file DOES NOT exist or
	// the user wants to OVERWRITE the existing file
	logger.Informational("Copying %s to %s\n", srcFile, destFile)
	err = ioutil.WriteFile(destFile, srcContents, 0644)
	if err != nil {
		logger.Error("There was an error trying to write to file '%s' (from original file '%s'). Exiting\n", destFile, srcContents)
		panic(err)
	}
}

func templateFile(srcContents []byte, templateVars config.Project, filename string) []byte {
	template, err := template.New(filename).Parse(string(srcContents))
	if err != nil {
		logger.Error("There was an error when parsing template from file '%s'. Exiting\n", filename)
		panic(err)
	}

	buf := new(bytes.Buffer)
	err = template.Execute(buf, templateVars)
	if err != nil {
		logger.Error("There was an error when executing template from file '%s'. Exiting\n", filename)
		panic(err)
	}

	return buf.Bytes()
}

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
