package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime"
)

type File struct {
	fileName         string
	similarFileNames []string
}

type Data struct {
	files []File
}

func CopyFile(file string) {
	_, currentDir, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("no caller information")
	}
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	srcFile := path.Join(path.Dir(currentDir), "files", file)
	destDir := path.Join(dir, file)
	if err != nil {
		fmt.Errorf("dir %v does not exists", err)
	}

	fmt.Printf("\033[38;5;205mcopying %s to %s\033[m\n", srcFile, destDir)

	// ensure directory exists
	os.MkdirAll(path.Dir(destDir), 0755)

	from, err := os.Open(srcFile)
	if err != nil {
		log.Fatal(err)
	}
	defer from.Close()

	to, err := os.OpenFile(destDir, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		log.Fatal(err)
	}
}

func GetFiles() Data {
	return Data{
		files: []File{
			{
				fileName:         ".editorconfig",
				similarFileNames: []string{},
			},
			{
				fileName:         ".gitattributes",
				similarFileNames: []string{},
			},
			{
				fileName:         "license",
				similarFileNames: []string{"LICENSE"},
			},
			{
				fileName:         ".github/CODE_OF_CONDUCT.md",
				similarFileNames: []string{},
			},
			{
				fileName:         ".github/COMMIT_CONVENTIONS.md",
				similarFileNames: []string{},
			},
			{
				fileName:         ".github/CONTRIBUTING.md",
				similarFileNames: []string{},
			},
			{
				fileName:         ".github/PULL_REQUEST_TEMPLATE.md",
				similarFileNames: []string{},
			},
			{
				fileName:         ".github/ISSUE_TEMPLATE/BUG_REPORT.md",
				similarFileNames: []string{},
			},
			{
				fileName:         ".github/ISSUE_TEMPLATE/config.yml",
				similarFileNames: []string{},
			},
			{
				fileName:         ".github/ISSUE_TEMPLATE/FEATURE_REQUEST.md",
				similarFileNames: []string{},
			},
			{
				fileName:         ".github/ISSUE_TEMPLATE/QUESTION.md",
				similarFileNames: []string{},
			},
		},
	}
}
