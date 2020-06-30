package main

func CopyOverExistingFiles() {
	for _, relativePath := range GetFilesToCopyOver() {
		path := GetFullPaths(relativePath)

		CopyFile(path.SrcPath, path.DestPath, relativePath)
	}
}

func main() {
	printInfo("Starting Globe!\n")

	CopyOverExistingFiles()
}
