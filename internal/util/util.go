package util

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"unsafe"

	"github.com/eankeen/dotty/internal/t"
	"github.com/eankeen/go-logger"
)

// HandleError panics if err is not nil
func HandleError(err error) {
	if err != nil {
		logger.Critical("%s\n", err)
		panic(err)
	}
}

func HandleFsError(err error) {
	if err == nil {
		return
	}

	if os.IsPermission(err) {
		logger.Error("You do not have permission to access the file or folder\n")
		log.Fatalln(err)
	}

	if os.IsNotExist(err) {
		logger.Error("File does not exist\n")
		log.Fatalln(err)
	}

	if os.IsExist(err) {
		logger.Error("File exists\n")
		log.Fatalln(err)
	}

	logger.Critical("An unknown error occurred\n")
	log.Panicln(err)
}

// Dirname performs same function as `__dirname()` in Node, obtaining the parent folder of the file of the callee of this function
func Dirname() string {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		panic("could not recover information from call stack")
	}

	dir := path.Dir(filename)
	return dir
}

// GetTtyWidth gets the tty's width, or number of columns
func GetTtyWidth() int {
	type winsize struct {
		Row    uint16
		Col    uint16
		Xpixel uint16
		Ypixel uint16
	}

	ws := &winsize{}
	retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		panic(errno)
	}
	return int(ws.Col)
}

// Contains tests to see if a particular string is in an array
func Contains(arr []string, str string) bool {
	for _, el := range arr {
		if el == str {
			return true
		}
	}
	return false
}

// Prompt ensures that we get a valid response
func Prompt(options []string, printText string, printArgs ...interface{}) string {
	logger.Notice(printText, printArgs...)

	var input string
	_, err := fmt.Scanln(&input)
	HandleError(err)

	if Contains(options, input) {
		return input
	}

	return Prompt(options, printText, printArgs)
}

func OpenPager(file string) {
	pager := os.Getenv("PAGER")
	program := "less"
	if pager != "" {
		program = pager
	}

	cmd := exec.Command(program, file)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	HandleError(err)
}

// PathExpand converts '~`, and to absolute path
func pathExpand(dotfilesDir string, rawPath string) string {
	expand := func(path string) string {
		homeDir, err := os.UserHomeDir()
		HandleError(err)

		newPath := strings.ReplaceAll(path, "~", homeDir)
		newPath = strings.ReplaceAll(path, "$HOME", homeDir)
		return newPath
	}

	isAbsolute := func(path string) bool {
		if strings.HasPrefix(path, "~") || strings.HasPrefix(path, "/") || strings.HasPrefix(path, "$HOME") {
			return true
		}

		return false
	}

	if isAbsolute(rawPath) {
		return expand(rawPath)
	}

	// relative
	return filepath.Join(dotfilesDir, expand(rawPath))
}

// Src gets
func Src(dotfilesDir string, dottyCfg t.DottyConfig, typ string) string {
	switch typ {
	case "system":
		if dottyCfg.SystemDirSrc == "" {
			return filepath.Join(dotfilesDir, "system")
		}
		return pathExpand(dotfilesDir, dottyCfg.SystemDirSrc)
	case "user":
		if dottyCfg.UserDirSrc == "" {
			return filepath.Join(dotfilesDir, "user")
		}
		return pathExpand(dotfilesDir, dottyCfg.UserDirSrc)
	case "local":
		if dottyCfg.LocalDirSrc == "" {
			return filepath.Join(dotfilesDir, "local")
		}
		return pathExpand(dotfilesDir, dottyCfg.LocalDirSrc)
	}

	panic("Src not valid")
}

// Dest is
func Dest(dotfilesDir string, dottyCfg t.DottyConfig, typ string) string {
	switch typ {
	case "system":
		if dottyCfg.SystemDirDest == "" {
			return "/"
		}
		return pathExpand(dotfilesDir, dottyCfg.SystemDirDest)
	case "user":
		if dottyCfg.UserDirDest == "" {
			homeDir, err := os.UserHomeDir()
			HandleFsError(err)

			return homeDir
		}
		return pathExpand(dotfilesDir, dottyCfg.UserDirDest)
	case "local":
		wd, err := os.Getwd()
		HandleError(err)

		return wd
	}

	panic("Dest not valid")
}
