package actions

import (
	"github.com/eankeen/dotty/fs"
)

// Apply (symlink) dotfiles
func Apply(dotfilesDir string, srcDir string, destDir string) {
	onFile := func(src string, dest string, rel string) {
		fs.ApplyFile(src, dest, rel)
	}

	onFolder := func(src string, dest string, rel string) {
		fs.ApplyFolder(src, dest, rel)
	}

	fs.Walk(dotfilesDir, srcDir, destDir, onFile, onFolder)
}

// Unapply (un-symlink) dotfiles
func Unapply(dotfilesDir string, srcDir string, destDir string) {
	onFile := func(src string, dest string, rel string) {
		fs.UnapplyFile(src, dest, rel)
	}

	onFolder := func(src string, dest string, rel string) {
		fs.UnapplyFolder(src, dest, rel)
	}

	fs.Walk(dotfilesDir, srcDir, destDir, onFile, onFolder)
}
