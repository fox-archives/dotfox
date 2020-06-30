package main

import "fmt"

func CopyOverExistingFiles() {
	for _, file := range getYaml().Files {
		printInfo("Processing file %s\n", file.Path)

		CopyFile(file)
		fmt.Println()
	}
}

func RemoveWrongFiles() {
	for _, file := range getYaml().OldFiles {
		printInfo("Removing bad file %s\n", file.Path)

		RemoveFile(file)
		fmt.Println()
	}
}

func main() {
	CopyOverExistingFiles()
	RemoveWrongFiles()
}
