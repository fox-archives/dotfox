import os
import parsetoml
import parseopt
import "./do"
import "./util"

proc writeHelp() =
  echo "help"

proc writeVersion() =
  echo "version"

if paramCount() < 1:
  echo "Error: Expected subcommand. Exiting"
  quit 1

var p = initOptParser(commandLineParams())
for kind, key, val in p.getopt():
  case kind
  of cmdEnd: break
  of cmdShortOption, cmdLongOption:
      case key
      of "help", "h": writeHelp(); quit 0
      of "version", "v": writeVersion(); quit 0
  of cmdArgument:
    let toml = parsetoml.parseFile(joinPath(getConfigDir(), "dotty/dotty.toml"))
    let dotDir = expandTilde(toml["config"]["dotDir"].getStr())
    let homeDir = expandTilde(toml["config"]["homeDir"].getStr())

    case key:
    of "status":
      doStatus(dotDir, homeDir, getDotFiles())
    of "reconcile":
      doReconcile(dotDir, homeDir, getDotFiles())
    else:
      echo "Error: Subcommand not recognized. Exiting"
      quit 1
