package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/eankeen/dotty/config"
	"github.com/eankeen/dotty/fs"
	"github.com/eankeen/dotty/internal/util"
	logger "github.com/eankeen/go-logger"
	"github.com/otiai10/copy"
	"github.com/spf13/cobra"
)

var userApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply updates intelligently",
	Run: func(cmd *cobra.Command, args []string) {
		storeDir := cmd.Flag("dot-dir").Value.String()
		destDir := cmd.Flag("dest-dir").Value.String()
		userDir := filepath.Join(storeDir, "user")

		userDotsConfig := config.GetUserToml(storeDir)

		err := filepath.Walk(userDir, func(src string, srcInfo os.FileInfo, err error) error {
			// prevent errors in slice
			if src == userDir {
				return nil
			}

			rel := src[len(userDir)+1:]
			dest := filepath.Join(destDir, rel)

			for _, file := range userDotsConfig.Files {
				logger.Debug("src: %s\n", src)
				logger.Debug("file.File: %s\n", file.File)

				if len(file.File) == 0 {
					logger.Warning("An entry in your '%s' config does not have a 'name' property. (This message may be repeated multiple times for each entry). Skipping\n", "user.dots.toml")
					return nil
				}

				fileMatches, fileType := config.FileMatches(src, file)

				if fileMatches && fileType == "file" {
					logger.Informational("Operating on File: '%s'\n", file.File)
					resolveFile(src, dest, rel)

					// go to next file (in dotfile folder)
					return nil
				} else if fileMatches && fileType == "folder" {
					logger.Informational("Operating on Folder: '%s'\n", file.File)
					resolveDirectory(src, dest, rel)

					// go to next file (in dotfile folder)
					return nil
				}

				// file in config did not match. continue because
				// another one might
				continue
			}

			// no explicit match was made. it may have already been matched, not match at all, or a parent folder matched instead
			return nil
		})
		util.P(err)
	},
}

func init() {
	userCmd.AddCommand(userApplyCmd)
}

// assumptions: src file exists, dest may NOT
func resolveFile(src string, dest string, rel string) {
	fi, err := os.Lstat(dest)
	if err != nil {
		// dest file doesn't exist
		if os.IsNotExist(err) {
			logger.Debug("OK: dest '%s' doesn't exist. Recreating\n", dest)
			err := fs.MkdirThenSymlink(src, dest)
			util.P(err)
			return
		}

		// some other issue (permissions, etc.)
		panic(err)
	}

	// dest file exists and is a symbolic link
	if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
		linkDest, err := os.Readlink(dest)
		util.P(err)

		logger.Debug("linkDest: %s\n", linkDest)
		logger.Debug("src: %s\n", src)
		// if link destination doesn't match src
		if linkDest != src {
			logger.Debug("OK: Symlink points to invalid location. Removing and Recreating\n")
			err := fs.RemoveThenSymlink(src, dest)
			util.P(err)
			return
		}

		// symlink does point to proper file,
		// no need to do anything
		return
	}

	// dest file exists and is NOT a symbolic link

	// read dest/src files to determine if they have the same content
	destContents, err := ioutil.ReadFile(dest)
	util.P(err)

	srcContents, err := ioutil.ReadFile(src)
	util.P(err)

	// if files have the same content
	if strings.Compare(string(destContents), string(srcContents)) == 0 {
		logger.Debug("OK: dest and src have same content. Replacing dest with symlink\n")
		err := fs.RemoveThenSymlink(src, dest)
		util.P(err)
		return
	}

	// files have different content
	// prompt user which to keep
promptUserFile:
	input := util.Prompt([]string{"compare", "use-src", "use-dest", "skip"}, "File %s exists both in your src and dest. (compare, use-src, use-dest, skip) ", rel)
	switch input {
	case "compare":
		// TODO: this could be cleaner
		var sep strings.Builder
		for i := 0; i < util.GetTtyWidth(); i++ {
			sep.WriteByte('-')
		}

		var output strings.Builder
		output.WriteString(fmt.Sprintf("SRC: %s\n", src))
		output.WriteString(sep.String())
		output.Write(srcContents)
		output.WriteString(sep.String())
		output.WriteString(fmt.Sprintf("\n\nDEST: %s\n", dest))
		output.WriteString(sep.String())
		output.Write(destContents)
		output.WriteString(sep.String())

		temp, err := ioutil.TempFile(os.TempDir(), "dotty-")
		defer os.Remove(temp.Name())
		util.P(err)

		_, err = temp.Write([]byte(output.String()))
		util.P(err)

		util.OpenEditor(temp.Name())

		goto promptUserFile
	case "use-src":
		err := fs.RemoveThenSymlink(src, dest)
		util.P(err)
		break
	case "use-dest":
		// copy file, replacing src
		destFile, err := os.Open(dest)
		util.P(err)
		defer destFile.Close()

		// umask will be applied after
		srcFile, err := os.Create(src)
		util.P(err)
		defer srcFile.Close()

		_, err = io.Copy(srcFile, destFile)
		util.P(err)

		// re-symlink
		err = fs.RemoveThenSymlink(src, dest)
		util.P(err)
		break
	case "skip":
		logger.Debug("Skipping '%s'\n", rel)
		break
	default:
		logger.Informational("Unknown Response\n")
	}

	return
}

// assumptions: src directory exists. dest may NOT
func resolveDirectory(src string, dest string, rel string) {
	fi, err := os.Lstat(dest)
	if err != nil {
		// dest file doesn't exist
		if os.IsNotExist(err) {
			logger.Debug("OK: dest '%s' doesn't exist. Recreating\n", dest)
			err := fs.MkdirThenSymlink(src, dest)
			util.P(err)
			return
		}

		// some other issue (permissions, etc.)
		panic(err)
	}

	// dest folder exists and is a symbolic link
	if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
		linkDest, err := os.Readlink(dest)
		util.P(err)

		logger.Debug("linkDest: %s\n", linkDest)
		logger.Debug("src: %s\n", src)
		// if link destination doesn't match src
		if linkDest != src {
			logger.Debug("OK: Symlink points to invalid location. Removing and Recreating\n")
			fs.RemoveThenSymlink(src, dest)
			util.P(err)
			return
		}

		// link has correct dest
		// no need to do anything
		return
	}

	// dest folder exists and is NOT a symbolic link

	srcDirs, err := ioutil.ReadDir(src)
	util.P(err)

	destDirs, err := ioutil.ReadDir(dest)
	util.P(err)

	// if both folders are empty, symlink them
	if len(srcDirs) == 0 && len(destDirs) == 0 {
		logger.Debug("OK: Replacing folder with symlink\n")
		err := fs.RemoveThenSymlink(src, dest)
		util.P(err)
	}

	if len(srcDirs) > 0 && len(destDirs) == 0 {
		logger.Debug("OK: Replacing folder with symlink\n")
		err := fs.RemoveThenSymlink(src, dest)
		util.P(err)
	}

	if len(destDirs) > 0 && len(srcDirs) == 0 {
		// copy the contents to source, and resymlink that
		err = copy.Copy(dest, src)
		util.P(err)

		err = os.RemoveAll(dest)
		util.P(err)

		err = os.Symlink(src, dest)
		util.P(err)
	}

	// Both srcDir and destDir have content
	// prompt user for which to keep
	if len(srcDirs) > 0 && len(destDirs) > 0 {
	promptUserFile:
		input := util.Prompt([]string{"compare", "use-src", "use-dest", "skip"}, "Folder %s exists both in your src and dest. (compare, use-src, use-dest, skip) ", rel)
		switch input {
		case "compare":
			fmt.Println("Not Implemented")
			goto promptUserFile
			// break
		case "use-src":
			fmt.Println("Not Implemented")
			goto promptUserFile
			// break
		case "use-dest":
			fmt.Println("Not Implemented")
			goto promptUserFile
			// break
		case "skip":
			logger.Debug("Skipping '%s'\n", rel)
			break
		}
	}

	return
}
