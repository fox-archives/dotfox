package e2e

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func _filename() string {
	_, _filename, _, ok := runtime.Caller(1)
	if !ok {
		log.Panicln("runtime.Caller not OK")
	}

	return _filename
}

func _dirname() string {
	return filepath.Dir(_filename())
}

func testDir() string {
	return filepath.Join(filepath.Dir(_dirname()), "testdata")
}

func ensureDir(t *testing.T, destDir string, path string) {
	finalDest := filepath.Join(destDir, path)
	_, err := os.Lstat(finalDest)
	if err != nil {
		if os.IsNotExist(err) {
			t.Logf("Error: ensureDir: Folder '%s' does not exist", finalDest)
		} else {
			t.Log("An unknown error occured")
			t.Error(err)
		}
	}

}

func ensureSymlink(t *testing.T, src string, dest string, path string) {
	finalSrc := filepath.Join(src, path)
	finalDest := filepath.Join(dest, path)

	linkSrc, err := os.Readlink(finalDest)
	if err != nil {
		// must not error
		t.Error(err)
		return
	}

	if linkSrc != finalSrc {
		t.Logf("Error: Symlink at '%s' points to '%s', but should point at '%s'", finalDest, linkSrc, finalSrc)
	}
}
