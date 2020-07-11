package globe

import (
	"fmt"

	"github.com/eankeen/globe/inspect"
	"github.com/eankeen/globe/internal/util"
)

// Update project with all bootstrap files
func Update(project inspect.Project) {
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
