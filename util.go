package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"

	"github.com/mattn/go-isatty"
	"gopkg.in/yaml.v2"
)

func FileExists(name string) (bool, error) {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err != nil, err
}

func getYaml() Files {
	yamlFile, err := ioutil.ReadFile(_dirname() + "/files.yml")
	if err != nil {
		printError("failed to read files.yml file")
		log.Print(err)
	}

	var y Files
	err = yaml.Unmarshal(yamlFile, &y)
	if err != nil {
		log.Fatalln(err)
	}

	return y
}

func _dirname() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatalln("could not recover information from call stack")
	}

	dir := path.Dir(filename)
	return dir
}

func isColorEnabled() bool {
	_, isNoColorEnabled := os.LookupEnv("NO_COLOR")
	if (os.Getenv("TERM") != "dumb") && !isNoColorEnabled &&
		isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		return true
	}
	return false
}

func debug(text string, args ...interface{}) {
	isDebug := func() bool {
		_, ok := os.LookupEnv("DEBUG")
		if !ok {
			return false
		}
		return true
	}

	if isDebug() {
		if isColorEnabled() {
			fmt.Print("\033[33m")
			fmt.Printf(text, args...)
			fmt.Print("\033[m")
			return
		}

		fmt.Printf(text, args...)
	}
}

func isColor() bool {
	_, isNoColorEnabled := os.LookupEnv("NO_COLOR")
	if (os.Getenv("TERM") != "dumb") && !isNoColorEnabled &&
		isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		return true
	}
	return false
}

func printInfo(text string, args ...interface{}) {
	if isColor() {
		fmt.Print("\033[32m")
		fmt.Printf(text, args...)
		fmt.Print("\033[m")
	} else {
		fmt.Printf(text, args...)
	}
}

func printError(text string, args ...interface{}) {
	if isColor() {
		fmt.Print("\033[31m")
		fmt.Printf(text, args...)
		fmt.Print("\033[m")
	} else {
		fmt.Printf(text, args...)
	}
}
