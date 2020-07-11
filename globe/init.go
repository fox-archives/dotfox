package globe

import (
	"os"
	"path"

	"github.com/eankeen/globe/internal/util"
)

func copyInitFiles(projectLocation string) {
	src := path.Join(util.Dirname(), "init")
	dest := path.Join(projectLocation)

	if err := copyDirRecurse(src, dest); err != nil {
		panic(err)
	}
}

// Init Globe config
func Init() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	copyInitFiles(wd)
}
