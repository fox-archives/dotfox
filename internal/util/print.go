package util

import (
	"fmt"
	"os"

	"github.com/mattn/go-isatty"
)

func isColorEnabled() bool {
	_, isNoColorEnabled := os.LookupEnv("NO_COLOR")
	if (os.Getenv("TERM") != "dumb") && !isNoColorEnabled &&
		isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		return true
	}
	return false
}

func isColor() bool {
	_, isNoColorEnabled := os.LookupEnv("NO_COLOR")
	if (os.Getenv("TERM") != "dumb") && !isNoColorEnabled &&
		isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		return true
	}
	return false
}

// PrintInfo passes params to `fmt.Printf`, colored as blue if it's a supporting tty
func PrintInfo(text string, args ...interface{}) {
	if isColor() {
		fmt.Print("\033[32m")
		fmt.Printf(text, args...)
		fmt.Print("\033[m")
	} else {
		fmt.Printf(text, args...)
	}
}

// PrintError passes params to `fmt.Printf`, colored as red if it's a supporting tty
func PrintError(text string, args ...interface{}) {
	if isColor() {
		fmt.Print("\033[31m")
		fmt.Printf(text, args...)
		fmt.Print("\033[m")
	} else {
		fmt.Printf(text, args...)
	}
}

// PrintDebug passes params to `fmt.Printf`, colored as blue if it's a supporting tty
func PrintDebug(text string, args ...interface{}) {
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
