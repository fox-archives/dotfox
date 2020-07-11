package inits

import (
	"os"
	"path"

	"github.com/eankeen/globe/internal/util"
)

func copyInitFiles(projectLocation string) {
	src := path.Join(util.Dirname(), "files")
	dest := path.Join(projectLocation)

	if err := copyDirRecurse(src, dest); err != nil {
		panic(err)
	}
}

// Inits Globe config
func Inits() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	copyInitFiles(wd)
}
