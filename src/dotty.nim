import os
import parsetoml
import parseopt
import strformat
import "./do"
import "./util"

var options = Options(showOk: true, isRoot: false, action: "")

var p = initOptParser(commandLineParams())
for kind, key, val in p.getopt():
  case kind
  of cmdShortOption, cmdLongOption:
    case key
    of "help", "h":
      writeHelp()
      quit QuitSuccess
    of "version", "v":
      writeVersion()
      quit QuitSuccess
    of "show-ok":
      options.showOk = false
    of "root":
      options.isRoot = true
  of cmdArgument:
    case key:
    of "status":
      options.action = "status"
    of "reconcile":
      options.action = "reconcile"
    else:
      die fmt"Subcommand '{key}' not recognized"
  of cmdEnd:
    break

let tomlFile = joinPath(getConfigDir(), "dotty/config.toml")
if not fileExists(tomlFile):
  die fmt"{tomlFile} not found"

let toml = parsetoml.parseFile(tomlFile)
let dotDir = expandTilde(toml["config"]["dotDir"].getStr())
let homeDir = expandTilde(toml["config"]["destDir"].getStr())

case options.action:
of "status":
  if options.isRoot:
    ensureRoot()
    ensureRootFileOwnership(dotDir)
    doStatus(dotDir, homeDir, options, getDotFiles("dottyRoot.sh"))
  else:
    ensureNotRoot()
    doStatus(dotDir, homeDir, options, getDotFiles("dotty.sh"))
of "reconcile":
  if options.isRoot:
    ensureRoot()
    ensureRootFileOwnership(dotDir)
    doReconcile(dotDir, homeDir, options, getDotFiles("dottyRoot.sh"))
  else:
    ensureNotRoot()
    doReconcile(dotDir, homeDir, options, getDotFiles("dotty.sh"))
else:
  logError "Expected subcommand"
  writeHelp()
  quit QuitFailure