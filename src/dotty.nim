import os
import parsetoml
import parseopt
import strformat
import posix
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
    of "files":
      options.files = parseCategories(val)
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
    echo fmt"Not all files in {dotDir} are owned by root. Fix this"
    quit QuitFailure

  if len(options.files) == 0:
    let cfg = joinPath(getConfigDir(), "dotty", "dottyRoot.sh")
    options.files = @[cfg]

else:
  if geteuid() == 0:
    die "Must NOT be running as root"

  if len(options.files) == 0:
    let cfg = joinPath(getConfigDir(), "dotty", "dotty.sh")
    options.files = @[cfg]

case options.action:
of "status":
  doStatus(dotDir, homeDir, options, getDotFiles(options.files))
of "reconcile":
  doReconcile(dotDir, homeDir, options, getDotFiles(options.files))
else:
  logError "Expected subcommand"
  writeHelp()
  quit QuitFailure