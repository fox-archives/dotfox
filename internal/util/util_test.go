package util

import (
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

func TestPathExpand(t *testing.T) {
	homeDir, _ := os.UserHomeDir()

	root := filepath.Join(homeDir, ".my-dotfiles")

	paths := []struct {
		Input    string
		Expected string
	}{
		{"~", os.Getenv("HOME")},
		{"$XDG_CONFIG_HOME/folder", filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "folder")},
		{"$HOME/folder2", filepath.Join(os.Getenv("HOME"), "folder2")},
		{"/", "/"},
		{"system", filepath.Join(root, "system")},
	}

	for _, path := range paths {
		actual := pathExpand(root, path.Input)

		if actual != path.Expected {
			t.Errorf("\nInput: '%s'\nExpected: '%s'\nActual: '%s'", path.Input, path.Expected, actual)
		}
	}
}
