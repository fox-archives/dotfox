package actions

import (
	"github.com/eankeen/dotty/fs"
)

// Apply (symlink) dotfiles
func Apply(dotfilesDir string, srcDir string, destDir string) {
	onFile := func(src string, dest string, rel string, mode int) {
		fs.ApplyFile(src, dest, rel, mode)
	}

	onFolder := func(src string, dest string, rel string, mode int) {
		fs.ApplyFolder(src, dest, rel, mode)
	}

	fs.Walk(dotfilesDir, srcDir, destDir, onFile, onFolder)
}

// Unapply (un-symlink) dotfiles
func Unapply(dotfilesDir string, srcDir string, destDir string) {
	onFile := func(src string, dest string, rel string, mode int) {
		fs.UnapplyFile(src, dest, rel)
	}

	onFolder := func(src string, dest string, rel string, mode int) {
		fs.UnapplyFolder(src, dest, rel)
	}

	fs.Walk(dotfilesDir, srcDir, destDir, onFile, onFolder)
}
