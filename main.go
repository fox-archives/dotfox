package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
)

func CcopyFile(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

func CopyDirRecurse(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = CopyDirRecurse(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = CcopyFile(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

func CopyInitFiles() {
	projectDir := getProjectDir()

	src := path.Join(_dirname(), "init")
	dest := path.Join(projectDir)
	err := CopyDirRecurse(src, dest)

	if err != nil {
		fmt.Println("Error when copying init files")
		panic(err)
	}
}

func flags() {
	boolPtr := flag.Bool("init", false, "Initiate a new Globe-managed project")
	flag.Parse()

	debug("initPtr: %v\n", *boolPtr)
	if *boolPtr == true {
		CopyInitFiles()
		printInfo("Initiated a new Globe-managed project\n")
		os.Exit(0)
	}
}

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
	flags()
	CopyOverExistingFiles()
	RemoveWrongFiles()
}
