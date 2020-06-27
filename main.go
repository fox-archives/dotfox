package main

import "os"

func main() {
	for _, relativePath := range GetFilesToCopy() {
		path := GetFullPaths(relativePath.name)

		CopyFile(path.SrcPath, path.DestPath, relativePath.name)
	}

	for _, relativePath := range GetFilesThatWereReplaced() {
		path := GetFullPaths(relativePath.name)
		// prompt to remove preexisting file if it exists
		stat, _ := os.Stat(path.DestPath)
		if stat != nil {
			// file exists, we ask if we should remove file
			shouldRemove := ShouldRemoveExistingFile(path.DestPath, relativePath.name)
			if shouldRemove == false {
				return
			}

			os.Remove(path.DestPath)
		}
	}
}
