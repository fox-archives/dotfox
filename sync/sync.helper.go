package sync

import (
	"bufio"
	"fmt"
	"os"

	"github.com/eankeen/globe/config"
	"github.com/eankeen/globe/internal/util"
	"github.com/gobwas/glob"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func projectFilesContain(project config.Project, glob glob.Glob) bool {
	files, err := util.GetChildFilesRecurse(project.StoreDir)
	if err != nil {
		panic(err)
	}

	var doesContain bool
	for _, file := range files {
		if glob.Match(file) {
			doesContain = true
			break
		}
	}

	util.PrintDebug("Does project contain pattern %+v?: %t\n", glob, doesContain)
	return doesContain
}

func shouldRemoveExistingFile(path string, relativePath string, destContents []byte, srcContents []byte) bool {
	util.PrintInfo("FileEntry '%s' is outdated. Replace it? (y/d/n): ", relativePath)
	r := bufio.NewReader(os.Stdin)
	c, err := r.ReadByte()
	if err != nil {
		panic(err)
	}

	if c == byte('Y') || c == byte('y') {
		util.PrintInfo("chose: yes\n")
		return true
	} else if c == byte('N') || c == byte('n') {
		util.PrintInfo("chose: no\n")
		return false
	} else if c == byte('D') || c == byte('d') {
		util.PrintInfo("chose: diff\n")
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(string(destContents), string(srcContents), true)
		fmt.Println(dmp.DiffPrettyText(diffs))
		return shouldRemoveExistingFile(path, relativePath, destContents, srcContents)
	} else {
		return shouldRemoveExistingFile(path, relativePath, destContents, srcContents)
	}
}

func isFileRelevant(project config.Project, file config.FileEntry) bool {
	projectContainsGoFiles := func() bool {
		if projectFilesContain(project, glob.MustCompile("*.go")) {
			return true
		}
		return false

	}

	switch file.For {
	case "all":
		return true
	case "golang":
		if projectContainsGoFiles() {
			return true
		}
		return false
	}

	util.PrintDebug("FileEntry '%s' does not match case statement. Has value %s. Skipping\n", file.RelPath, file.For)
	return false
}
