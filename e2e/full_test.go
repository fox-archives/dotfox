package e2e

import (
	"log"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/eankeen/dotty/fs"
)

func run() {

}

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

func do(dotDir string, srcDir string, destDir string) {
	onFile := func(src string, dest string, rel string) {
		fs.ApplyFile(src, dest, rel)
	}

	onFolder := func(src string, dest string, rel string) {
		fs.ApplyFolder(src, dest, rel)
	}

	fs.Walk(dotDir, srcDir, destDir, onFile, onFolder)

	time.Sleep(time.Millisecond * 500)

	// unlink
	onFile2 := func(src string, dest string, rel string) {
		fs.UnapplyFile(src, dest, rel)
	}

	onFolder2 := func(src string, dest string, rel string) {
		fs.UnapplyFolder(src, dest, rel)
	}

	fs.Walk(dotDir, srcDir, destDir, onFile2, onFolder2)
}

func TestFull(t *testing.T) {
	testDir := filepath.Join(filepath.Dir(_dirname()), "testdata")
	test1 := filepath.Join(testDir, "test1")

	dotDir := test1
	srcDir := filepath.Join(test1, "dotfiles")
	destDir := filepath.Join(test1, "user-home")
	do(dotDir, srcDir, destDir)

	t.Log("thing")
}
