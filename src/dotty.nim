import os
import parsetoml
import parseopt
import strformat
import "./do"
import "./util"

var options = Options(showOk: true, isRoot: false, interactive: false)

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
    of "interactive", "i":
      options.interactive = true
    of "show-ok":
      options.showOk = false
    of "config":
      options.configFile = val
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

if options.configFile == "":
  options.configFile = joinPath(getConfigDir(), "dotty", "config.toml")

if not fileExists(options.configFile):
  die fmt"Config file '{options.configFile}' not found"

let toml = parsetoml.parseFile(options.configFile)
let dotDir = expandTilde(toml["config"]["dotDir"].getStr())
let homeDir = expandTilde(toml["config"]["destDir"].getStr())

var scriptName = ""
if options.isRoot:
  ensureRoot()
  ensureRootFileOwnership(dotDir)
  scriptName = "dotty.sh"
else:
  ensureNotRoot()
  scriptName = "dotty.sh"

case options.action:
of "status":
  doStatus(dotDir, homeDir, options, getDotFiles(scriptName))
of "reconcile":
  doReconcile(dotDir, homeDir, options, getDotFiles(scriptName))
else:
  logError "Expected subcommand"
  writeHelp()
  quit QuitFailure