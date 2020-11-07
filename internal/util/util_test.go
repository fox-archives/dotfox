package util

import (
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestContains(t *testing.T) {
	arr := []string{"alfa", "bravo"}

	// test string
	if !Contains(arr, "alfa") {
		t.Error("string 'alfa' not in array 'arr'. it is supposed to be in it")
	}
}

func TestCreatePath(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Panicln(homeDir)
	}

	root := filepath.Join(homeDir, ".my-dotfiles")

	paths := []struct {
		Input  string
		Output string
	}{
		{"~", "/home/edwin"},
		{"/", "/"},
		{"system", filepath.Join(root, "system")},
	}

	for _, path := range paths {
		newPath := CreatePath(root, path.Input)

		if newPath != path.Output {
			t.Errorf("'%s' is not '%s'", newPath)
		}
	}
}
