import os
import parsetoml
import parseopt
import strformat
import posix
import "./do"
import "./util"

var options = Options(showOk: true, isRoot: false)

var p = initOptParser(commandLineParams())
for kind, key, val in p.getopt():
  case kind
  of cmdShortOption, cmdLongOption:
    case key
    of "help", "h":
      writeHelp()
      quit QuitSuccess
    of "version", "v":
      echo "v0.5.0"
      quit QuitSuccess
    of "show-ok":
      options.showOk = parseBoolFlag(val)
    of "config":
        options.configFile = val
    of "root":
      options.isRoot = parseBoolFlag(val)
    of "deployment":
      options.deployment = val
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

if options.isRoot:
  if geteuid() != 0:
    die "Must be running as root"

  if not hasAllRootFiles(dotDir):
    die fmt"Not all files in {dotDir} are owned by root. Fix this"

  if options.deployment == "":
    options.deployment = joinPath(getConfigDir(), "dotty", "dottyRoot.sh")
else:
  if geteuid() == 0:
    die "Must NOT be running as root"

  if options.deployment == "":
    options.deployment = joinPath(getConfigDir(), "dotty", "dotty.sh")

case options.action:
of "status":
  doStatus(dotDir, homeDir, options, getDotfileList(options.deployment))
of "reconcile":
  doReconcile(dotDir, homeDir, options, getDotfileList(options.deployment))
else:
  logError "Expected subcommand"
  writeHelp()
  quit QuitFailure