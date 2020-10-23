package fs

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/eankeen/dotty/internal/util"
	logger "github.com/eankeen/go-logger"
	"github.com/otiai10/copy"
)

// assumptions: src file exists, dest may NOT
func ApplyFile(src string, dest string, rel string) {
	fi, err := os.Lstat(dest)
	if err != nil {
		// dest file doesn't exist
		if os.IsNotExist(err) {
			logger.Debug("OK: dest '%s' doesn't exist. Recreating\n", dest)
			err := MkdirThenSymlink(src, dest)
			// if err != nil {
			// 	panic(err)
			// }
			util.HandleFsError(err)
			return
		}

		// some other issue (permissions, etc.)
		panic(err)
	}

	// dest file exists and is a symbolic link
	if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
		linkDest, err := os.Readlink(dest)
		util.HandleFsError(err)

		logger.Debug("linkDest: %s\n", linkDest)
		logger.Debug("src: %s\n", src)
		// if link destination doesn't match src
		if linkDest != src {
			logger.Debug("OK: Symlink points to invalid location. Removing and Recreating\n")
			err := RemoveThenSymlink(src, dest)
			util.HandleFsError(err)
			return
		}

		// symlink does point to proper file,
		// no need to do anything
		return
	}

	// dest file exists and is NOT a symbolic link

	// read dest/src files to determine if they have the same content
	destContents, err := ioutil.ReadFile(dest)
	util.HandleFsError(err)

	srcContents, err := ioutil.ReadFile(src)
	util.HandleFsError(err)

	// if files have the same content
	if strings.Compare(string(destContents), string(srcContents)) == 0 {
		logger.Debug("OK: dest and src have same content. Replacing dest with symlink\n")
		err := RemoveThenSymlink(src, dest)
		util.HandleFsError(err)
		return
	}

	// files have different content
	// prompt user which to keep
promptUserFile:
	input := util.Prompt([]string{"diff", "diff2", "use-src", "use-dest", "skip"}, "File %s exists both in your src and dest. (diff, diff2, use-src, use-dest, skip) ", rel)
	switch input {
	case "diff":
		cmd := exec.Command("colordiff", src, dest)

		output, err := cmd.Output()
		if err != nil && err.Error() != "exit status 1" {
			panic(err)
		}

		file := WriteTemp(output)
		util.OpenPager(file.Name())

		goto promptUserFile
	case "diff2":
		cmd := exec.Command("colordiff", "--side-by-side", src, dest)

		output, err := cmd.Output()
		if err != nil && err.Error() != "exit status 1" {
			panic(err)
		}

		file := WriteTemp(output)
		util.OpenPager(file.Name())

		goto promptUserFile
	case "use-src":
		err := RemoveThenSymlink(src, dest)
		util.HandleFsError(err)
		break
	case "use-dest":
		// copy file, replacing src
		destFile, err := os.Open(dest)
		defer destFile.Close()
		util.HandleFsError(err)

		// umask will be applied after
		srcFile, err := os.Create(src)
		defer srcFile.Close()
		util.HandleFsError(err)

		_, err = io.Copy(srcFile, destFile)
		util.HandleFsError(err)

		// re-symlink
		err = RemoveThenSymlink(src, dest)
		util.HandleFsError(err)
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
func ApplyDirectory(src string, dest string, rel string) {
	fi, err := os.Lstat(dest)
	if err != nil {
		// dest file doesn't exist
		if os.IsNotExist(err) {
			logger.Debug("OK: dest '%s' doesn't exist. Recreating\n", dest)
			err := MkdirThenSymlink(src, dest)
			util.HandleFsError(err)
			return
		}

		// some other issue (permissions, etc.)
		panic(err)
	}

	// dest folder exists and is a symbolic link
	if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
		linkDest, err := os.Readlink(dest)
		util.HandleFsError(err)

		logger.Debug("linkDest: %s\n", linkDest)
		logger.Debug("src: %s\n", src)
		// if link destination doesn't match src
		if linkDest != src {
			logger.Debug("OK: Symlink points to invalid location. Removing and Recreating\n")
			RemoveThenSymlink(src, dest)
			util.HandleFsError(err)
			return
		}

		// link has correct dest
		// no need to do anything
		return
	}

	// dest folder exists and is NOT a symbolic link

	srcDirs, err := ioutil.ReadDir(src)
	util.HandleFsError(err)

	destDirs, err := ioutil.ReadDir(dest)
	util.HandleFsError(err)

	// if both folders are empty, symlink them
	if len(srcDirs) == 0 && len(destDirs) == 0 {
		logger.Debug("OK: Replacing folder with symlink\n")
		err := RemoveThenSymlink(src, dest)
		util.HandleFsError(err)
		return
	}

	if len(srcDirs) > 0 && len(destDirs) == 0 {
		logger.Debug("OK: Replacing folder with symlink\n")
		err := RemoveThenSymlink(src, dest)
		util.HandleFsError(err)
		return
	}

	if len(destDirs) > 0 && len(srcDirs) == 0 {
		// copy the contents to source, and re-symlink that
		err = copy.Copy(dest, src)
		util.HandleFsError(err)

		err = os.RemoveAll(dest)
		util.HandleFsError(err)

		err = os.Symlink(src, dest)
		util.HandleFsError(err)
		return
	}

	// Both srcDir and destDir have content
	// prompt user for which to keep
	if len(srcDirs) > 0 && len(destDirs) > 0 {
	promptUserFolder:
		input := util.Prompt([]string{"diff", "use-src", "use-dest", "skip"}, "Folder %s exists both in your src and dest. (diff, use-src, use-dest, skip) ", rel)
		switch input {
		case "diff":
			cmd := exec.Command("tree", src)

			output, err := cmd.Output()
			util.HandleError(err)

			cmd2 := exec.Command("tree", dest)
			output2, err2 := cmd2.Output()
			util.HandleError(err2)

			content := append(output, "\n\n"...)
			content = append(content, output2...)

			file := WriteTemp(content)
			util.OpenPager(file.Name())

			goto promptUserFolder
			// break
		case "use-src":
			err := os.RemoveAll(dest)
			util.HandleFsError(err)

			err = os.Symlink(src, dest)
			util.HandleFsError(err)

			break
		case "use-dest":
			// TODO: fix this (getting permissions errors)
			err := os.RemoveAll(src)
			util.HandleFsError(err)

			err = copy.Copy(dest, src)
			util.HandleFsError(err)

			err = os.RemoveAll(dest)
			util.HandleFsError(err)

			err = os.Symlink(src, dest)
			util.HandleFsError(err)

			break
		case "skip":
			logger.Debug("Skipping '%s'\n", rel)
			break
		}
	}
}

func UnapplyDirectory(src string, dest string, rel string) {
	fi, err := os.Lstat(dest)
	if os.IsNotExist(err) {
		return
	}
	util.HandleFsError(err)

	if fi.Mode()&os.ModeSymlink != os.ModeSymlink {
		logger.Error("Skipping: Not a Symlink: '%s'\n", dest)
		return
	}

	cmd := exec.Command("unlink", dest)
	res, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(res)
		util.HandleError(err)
	}

	// if fi.IsDir() {
	// 	return
	// }

	// realDest, err := os.Readlink(dest)
	// util.HandleFsError(err)

	// fmt.Println(realDest)
	// err = syscall.Unlink(realDest)
	// util.HandleError(err)
}
