package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/eankeen/globe/inits"
	"github.com/eankeen/globe/internal/util"
	"github.com/eankeen/globe/scan"
	"github.com/eankeen/globe/sync"
)

// init subcommand
func startInits() {
	inits.Inits()
}

// run subcommand
func startSync() {
	project := scan.Scan()

	util.PrintInfo("Project located at %s\n", project.ProjectLocation)
	sync.Sync(project)
}

func showHelp() {
	const s = `Command:
  globe

Description:
  An easy to use language-agnostic configuration management tool

Commands:
  init     Initiate Globe configuration
  sync   Update configuration and files

Options:
  --help   Display help menu
`
	fmt.Print(s)
}

// Run initiates the CLI
func Run() {
	util.PrintDebug("args: %v\n", os.Args)

	syncCmd := flag.NewFlagSet("sync", flag.ExitOnError)
	syncCmdOverride := syncCmd.Bool("override", false, "Overrides all existing files. Good for a non-interactive mode")

	if len(os.Args) == 1 {
		util.PrintError("No subcommand found. See `globe --help`\n")
		os.Exit(1)
		return
	}

	switch os.Args[1] {
	case "init":
		startInits()
	case "sync":
		syncCmd.Parse(os.Args[2:])
		if *syncCmdOverride == true {
			startSync()
		} else {
			startSync()
		}
	case "--help", "-help", "-h":
		showHelp()

	default:
		util.PrintError("Invalid arguments. See `globe --help`\n")
		os.Exit(1)
	}
}
