package sync

import (
	"fmt"

	"github.com/eankeen/globe/internal/util"
	"github.com/eankeen/globe/scan"
)

// Sync project with all bootstrap files
func Sync(project scan.Project) {
	// NEW FILES
	{
		for _, file := range project.BootstrapFiles.NewFiles {
			util.PrintInfo("Processing file %s\n", file.RelPath)

			copyFile(project, file)
			fmt.Println()
		}
	}

	// OLD FILES
	{
		for _, file := range project.BootstrapFiles.OldFiles {
			util.PrintInfo("Removing bad file %s\n", file.RelPath)

			removeFile(project, file)
			fmt.Println()
		}
	}
}
