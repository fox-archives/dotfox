package actions

import (
	"os"
	"os/exec"

	"github.com/eankeen/dotty/internal/util"
)

// OpenEditor opens a file for editing
func OpenEditor(file string) {
	editor := os.Getenv("EDITOR")
	program := "vim"
	if editor != "" {
		program = editor
	}

	cmd := exec.Command(program, file)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	util.HandleError(err)
}
