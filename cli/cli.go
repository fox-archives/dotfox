package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/eankeen/globe/globe"
	"github.com/eankeen/globe/inspect"
	"github.com/eankeen/globe/internal/util"
)

// init subcommand
func initStart() {
	globe.Init()
}

// run subcommand
func updateStart() {
	project := inspect.Inspect()

	util.PrintInfo("Project located at %s\n", project.ProjectLocation)
	globe.Update(project)
}

func showHelp() {
	const s = `Command:
  globe

Description:
  An easy to use language-agnostic configuration management tool

Commands:
  init     Initiate Globe configuration
  update   Update configuration and files

Options:
  --help   Display help menu
`
	fmt.Print(s)
}

// Run initiates the CLI
func Run() {
	util.PrintDebug("args: %v\n", os.Args)

	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
	updateCmdOverride := updateCmd.Bool("override", false, "Overrides all existing files. Good for a non-interactive mode")

	if len(os.Args) == 1 {
		util.PrintError("No subcommand found. See `globe --help`\n")
		os.Exit(1)
		return
	}

	switch os.Args[1] {
	case "init":
		initStart()
	case "update":
		updateCmd.Parse(os.Args[2:])
		if *updateCmdOverride == true {
			updateStart()
		} else {
			updateStart()
		}
	case "--help", "-help", "-h":
		showHelp()

	default:
		util.PrintError("Invalid arguments. See `globe --help`\n")
		os.Exit(1)
	}
}
