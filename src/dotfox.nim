import os
import parsetoml
import parseopt
import strformat
import posix
import "./do"
import "./util"

var options = Options(showOk: false, isRoot: false)

var p = initOptParser(commandLineParams())
for kind, key, val in p.getopt():
  case kind
  of cmdShortOption, cmdLongOption:
    case key
    of "help", "h":
      writeHelp()
      quit QuitSuccess
    of "version", "v":
      echo "v0.6.1"
      quit QuitSuccess
    of "show-ok":
      options.showOk = parseBoolFlag(val)
    of "config-dir":
      options.configDir = val
    of "root":
      options.isRoot = parseBoolFlag(val)
    of "deployment":
      options.deployment = val
    else:
      die fmt"Flag '{key}' not recognized"
  of cmdArgument:
    options.action = key
  of cmdEnd:
    break


if options.configDir == "":
  options.configDir = joinPath(getConfigDir(), "dotfox")

if not isAbsolute(options.configDir):
  die fmt"Directory '{options.configDir}' is not an absolute path"
if not dirExists(options.configDir):
  die fmt"Config directory '{options.configDir}' is not a directory"

let toml = parsetoml.parseFile(joinPath(options.configDir, "config.toml"))
let dotDir = expandTilde(toml["config"]["dotDir"].getStr())
let homeDir = expandTilde(toml["config"]["destDir"].getStr())

if options.isRoot:
  if geteuid() != 0:
    die "Must be running as root"

  if not hasAllRootFiles(dotDir):
    die fmt"Not all files in '{dotDir}' are owned by root. Fix this"
else:
  if geteuid() == 0:
    die "Must NOT be running as root"

if options.deployment == "":
  die "No deployment specified"

if options.action == "":
  die "Must pass a subcommand"

case options.action:
of "status":
  doStatus(dotDir, homeDir, options, getDotfileList(options))
of "deploy":
  doDeploy(dotDir, homeDir, options, getDotfileList(options))
of "debug":
  doDebug(dotDir, homeDir, options, getDotfileList(options))
else:
  logError fmt"Subcommand '{options.action}' not recognized"
  writeHelp()
  quit QuitFailure
