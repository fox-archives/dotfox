package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
)

type File struct {
	fileName         string
	similarFileNames []string
}

type Data struct {
	files []File
}

func CopyFile(src string, dest string) {
	dir, err := os.Getwd()
	fileDir := path.Join(dir, "files")
	destDir := path.Join(dir, dest)

	if err != nil {
		fmt.Errorf("dir %v does not exists", err)
	}
	fmt.Println(dir)

	from, err := os.Open(path.Join(fileDir, src))
	if err != nil {
		log.Fatal(err)
	}
	defer from.Close()

	to, err := os.OpenFile(destDir, os.O_RDWR|os.O_CREATE, 0666)
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
				fileName:         ".github/code_of_conduct.md",
				similarFileNames: []string{".github/CODE_OF_CONDUCT.md"},
			},
			{
				fileName:         ".github/commit_conventions.md",
				similarFileNames: []string{".github/COMMIT_CONVENTIONS.md"},
			},
			{
				fileName:         ".github/contributing.md",
				similarFileNames: []string{".github/CONTRIBUTING.md"},
			},
			{
				fileName:         ".github/pull_request_template.md",
				similarFileNames: []string{".github/PULL_REQUEST_TEMPLATE.md"},
			},
			{
				fileName:         ".github/issue_template/bug_report.md",
				similarFileNames: []string{".github/issue_template/BUG_REPORT.md", ".github/ISSUE_TEMPLATE/BUG_REPORT.md", ".github/ISSUE_TEMPLATE/bug_report.md"},
			},
			{
				fileName:         ".github/issue_template/config.yml",
				similarFileNames: []string{".github/issue_template/BUG_REPORT.md", ".github/ISSUE_TEMPLATE/BUG_REPORT.md", ".github/ISSUE_TEMPLATE/bug_report.md"},
			},
			{
				fileName:         ".github/issue_template/feature_request",
				similarFileNames: []string{".github/issue_template/FEATURE_REQUEST.md", ".github/ISSUE_TEMPLATE/FEATURE_REQUEST.md", ".github/ISSUE_TEMPLATE/feature_request.md"},
			},
			{
				fileName:         ".github/issue_template/question",
				similarFileNames: []string{".github/issue_template/QUESTION.md", ".github/ISSUE_TEMPLATE/QUESTION.md", ".github/ISSUE_TEMPLATE/question.md"},
			},
		},
	}
}
