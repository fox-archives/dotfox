package main

import (
	"log"
	"os"
	"path"
)

func main() {
	data := GetFiles()

	for i := 0; i < len(data.files); i++ {
		file := data.files[i]

		// first remove all similarFileNames
		for _, similarFileName := range file.similarFileNames {
			dir, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}
			fullPath := path.Join(dir, similarFileName)
			// fmt.Printf("removing %v\n", fullPath)
			errR := os.Remove(fullPath)
			if errR != nil {
				// file doesn't exist. that's okay
				// fmt.Printf(err)
			}
		}

		CopyFile(".editorconfig", "outputfile")
	}

}
