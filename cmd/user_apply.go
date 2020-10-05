package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/eankeen/globe/config"
	"github.com/eankeen/globe/internal/util"
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
			err := config.CreateNewSymlink(src, dest)
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
			err := config.FixBrokenSymlink(src, dest)
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
		err := os.Remove(dest)
		util.P(err)

		err = os.Symlink(src, dest)
		util.P(err)
		return
	}

	// files have different content
	// intelligently prompt user which to keep
promptUser:
	input := Prompt([]string{"compare", "use-src", "use-dest", "skip"}, "File %s exists both in your src and dest. (compare, use-src, use-dest, skip) ", rel)
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

		fmt.Println(string(output.String()))

		// TODO: display in a pager
		// cmd := exec.Command("sh", "-c", "echo '"+output.String()+"' | less")
		// cmd.Stdin = os.Stdin
		// cmd.Stdout = os.Stdout
		// cmd.Stderr = os.Stderr
		// stderr, _ := cmd.StderrPipe()

		// err = cmd.Start()
		// // data, err := cmd.CombinedOutput()
		// if err != nil {
		// 	log.Println(err)
		// }
		// // fmt.Println(string(data))
		// scanner := bufio.NewScanner(stderr)
		// for scanner.Scan() {
		// 	fmt.Println(scanner.Text())
		// }

		// fmt.Println(output.String())

		goto promptUser
	case "use-src":
		err := os.Remove(dest)
		util.P(err)

		err = os.Symlink(src, dest)
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
		err = os.Remove(dest)
		util.P(err)

		err = os.Symlink(src, dest)
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

// Prompt ensures that we get a valid response
func Prompt(options []string, printText string, printArgs ...interface{}) string {
	logger.Informational(printText, printArgs...)

	var input string
	_, err := fmt.Scanln(&input)
	util.P(err)

	if util.Contains(options, input) {
		return input
	}

	return Prompt(options, printText, printArgs)
}

// assumptions: src directory exists. dest may NOT
func resolveDirectory(src string, dest string, rel string) {
	fi, err := os.Lstat(dest)
	if err != nil {
		// dest file doesn't exist
		if os.IsNotExist(err) {
			err := config.CreateNewSymlink(src, dest)
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
			config.FixBrokenSymlink(src, dest)
			util.P(err)
			return
		}

		// link has correct dest
		// no need to do anything
		return
	}

	// dest folder exists and is NOT a symbolic link

	// make sure the src folder is empty before we copy it
	srcDirs, err := ioutil.ReadDir(src)
	util.P(err)

	// srcDir has content. we aren't sure if we want to override it
	// with content in destDir
	// TODO: intelligently prompt user
	if len(srcDirs) != 0 {
		logger.Warning("Skipping '%s' because the directory has content at %s. Remove the contents, or remove the directory at dest if you don't want that to take precedence", rel, src)
	}

	// copy the contents to source, and resymlink that
	err = copy.Copy(dest, src)
	// TODO: if src already exists, will it error?
	util.P(err)

	err = os.RemoveAll(dest)
	util.P(err)

	err = os.Symlink(src, dest)
	util.P(err)

	return
}
