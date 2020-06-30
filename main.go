package main

import "fmt"

func CopyOverExistingFiles() {
	for _, file := range getYaml().Files {
		printInfo("Processing file %s\n", file.Path)

		CopyFile(file)
		fmt.Println()
	}
}

func main() {
	printInfo("Starting Globe!\n")

	CopyOverExistingFiles()
}
