package fs

import (
	"bufio"
	"fmt"
	"os"

	"github.com/eankeen/globe/internal/util"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func shouldRemoveExistingFile(destFile string, relativePath string, destContents []byte, srcContents []byte) bool {
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
		return shouldRemoveExistingFile(destFile, relativePath, destContents, srcContents)
	} else {
		return shouldRemoveExistingFile(destFile, relativePath, destContents, srcContents)
	}
}
