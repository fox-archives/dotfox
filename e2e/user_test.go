package e2e

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/eankeen/dotty/actions"
	"github.com/eankeen/dotty/config"
	"github.com/eankeen/dotty/internal/util"
)

func TestUserApply(t *testing.T) {
	type UserTest struct {
		dir string
		fn  func(*testing.T, string, string)
	}

	// tests; path relative to destDir
	userTests := []UserTest{
		// basic; single file
		{
			dir: "user-1",
			fn: func(t *testing.T, srcDir, destDir string) {
				ensureDir(t, destDir, "")
				ensureSymlink(t, srcDir, destDir, "bar")
			},
		},
		// deeply nested subdirectories
		{
			dir: "user-2",
			fn: func(t *testing.T, srcDir, destDir string) {
				ensureDir(t, destDir, "")
				ensureSymlink(t, srcDir, destDir, "one")
				ensureSymlink(t, srcDir, destDir, "subdir-1/two")
				ensureSymlink(t, srcDir, destDir, "subdir-2/subdir-2-1/three")
			},
		},
	}

	for _, userTest := range userTests {
		dotfilesDir := filepath.Join(testDir(), "test-user", userTest.dir)
		dottyCfg := config.DottyCfg(dotfilesDir)

		srcDir := util.Src(dotfilesDir, dottyCfg, "user")
		destDir := util.Dest(dotfilesDir, dottyCfg, "user")

		err := os.RemoveAll(destDir)
		util.HandleFsError(err)
		// time.Sleep(time.Millisecond * 1000)
		actions.Apply(dotfilesDir, srcDir, destDir)

		userTest.fn(t, srcDir, destDir)
	}

}
