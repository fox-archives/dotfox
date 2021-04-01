import os
import parsetoml
import parseopt
import "./do"
import "./util"

proc writeHelp() =
  echo """Dotty

Usage: dotty [subcommand] [flags]

Subcommands:
  status
    Views the status of user dotfiles
  reconcile
    Symlinks dotfiles to proper location and attempts to autofix mismatches
  rootStatus
    Views the status of root dotfiles
  rootReconcile
    Copies over dotfiles to /root directory

Flags:
  --help
  --version"""

proc writeVersion() =
  # TODO
  echo "0.2.1"

if paramCount() < 1:
  die "Expected subcommand"

var p = initOptParser(commandLineParams())
for kind, key, val in p.getopt():
  case kind
  of cmdEnd: break
  of cmdShortOption, cmdLongOption:
      case key
      of "help", "h": writeHelp(); quit QuitSuccess
      of "version", "v": writeVersion(); quit QuitSuccess
  of cmdArgument:
    let toml = parsetoml.parseFile(joinPath(getConfigDir(), "dotty/config.toml"))
    let dotDir = expandTilde(toml["config"]["dotDir"].getStr())
    let homeDir = expandTilde(toml["config"]["destDir"].getStr())

    case key:
    of "status":
      doStatus(dotDir, homeDir, getDotFiles())
    of "rootStatus":
      doRootStatus(dotDir, homeDir, getDotFiles())
    of "reconcile":
      doReconcile(dotDir, homeDir, getDotFiles())
    of "rootReconcile":
      doRootReconcile(dotDir, homeDir, getRootDotFiles())
    else:
      die "Subcommand not recognized"
