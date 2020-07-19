package logic

import (
	"github.com/eankeen/globe/config"
	"github.com/eankeen/globe/fs"
	"github.com/eankeen/globe/internal/util"
	"github.com/gobwas/glob"
)

// ProcessFiles processes each file according to their 'Op' operation specified by the user
func ProcessFiles(project config.Project, files []config.FileEntry) {
	for _, file := range files {
		util.PrintDebug("Processing file %s\n", file.RelPath)

		if file.Op == "add" {

			// validate to see if we should even be trying to copy the file
			// over. for example scripts/go.sh should only be copied when
			// there are .go files in the repository
			isFileRelevant := isFileRelevant(project.ProjectDir, file)
			if !isFileRelevant {
				util.PrintInfo("Skipping irrelevant file '%s'\n", file.RelPath)
				return
			}
			fs.CopyFile(file.SrcPath, file.DestPath, file.RelPath)
			continue
		} else if file.Op == "remove" {
			fs.RemoveFile(file.DestPath)
			continue
		}

		util.PrintError("File '%s's operation could not be read. Exiting.\n", file.RelPath)
	}
}

func isFileRelevant(projectDir string, file config.FileEntry) bool {
	projectContainsGoFiles := func() bool {
		files, err := util.GetChildFilesRecurse(projectDir)
		if err != nil {
			panic(err)
		}
		if projectFilesContain(files, glob.MustCompile("*.go")) {
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

func projectFilesContain(files []string, glob glob.Glob) bool {
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
