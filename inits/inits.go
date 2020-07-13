package inits

import (
	"os"
	"path"
)

// Inits Globe config
func Inits(storeDir string) {
	projectDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	{
		src := path.Join(storeDir, "init")
		dest := path.Join(projectDir)

		if err := copyDirRecurse(src, dest); err != nil {
			panic(err)
		}
	}
}
