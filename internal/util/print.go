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
		fmt.Print("\033[1;32m")
		fmt.Print("INFO : ▶ ")
		fmt.Print("\033[0;32m")
		fmt.Printf(text, args...)
		fmt.Print("\033[m")
	} else {
		fmt.Printf(text, args...)
	}
}

// PrintWarning passes params to `fmt.Printf`, colored as yellow if it's a supporting tty
func PrintWarning(text string, args ...interface{}) {
	if isColor() {
		fmt.Print("\033[1;33m")
		fmt.Print("WARNG: ▶ ")
		fmt.Print("\033[0;33m")
		fmt.Printf(text, args...)
		fmt.Print("\033[m")
	} else {
		fmt.Printf(text, args...)
	}
}

// PrintError passes params to `fmt.Printf`, colored as red if it's a supporting tty
func PrintError(text string, args ...interface{}) {
	if isColor() {
		fmt.Print("\033[1;36m")
		fmt.Print("ERROR: ▶ ")
		fmt.Print("\033[0;36m")
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
			fmt.Print("\033[1;36m")
			fmt.Print("DEBUG: ▶ ")
			fmt.Print("\033[0;36m")
			fmt.Printf(text, args...)
			fmt.Print("\033[m")
			return
		}

		fmt.Printf(text, args...)
	}
}
