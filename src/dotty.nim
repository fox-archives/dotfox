import os
import parsetoml
import parseopt
import strformat
import "./do"
import "./util"

if paramCount() < 1:
  die "Expected subcommand"

var options = Options(showOk: true)

var p = initOptParser(commandLineParams())
for kind, key, val in p.getopt():
  case kind
  of cmdEnd: break
  of cmdShortOption, cmdLongOption:
    case key
    of "help", "h": writeHelp(); quit QuitSuccess
    of "version", "v": writeVersion(); quit QuitSuccess
    of "show-ok":
      options.showOk = false
  of cmdArgument:
    let tomlFile = joinPath(getConfigDir(), "dotty/config.toml")
    if not fileExists(tomlFile):
      die fmt"{tomlFile} not found"

    let toml = parsetoml.parseFile(tomlFile)
    let dotDir = expandTilde(toml["config"]["dotDir"].getStr())
    let homeDir = expandTilde(toml["config"]["destDir"].getStr())

    case key:
    of "status":
      ensureNotRoot()
      doStatus(dotDir, homeDir, options, getDotFiles("dotty.sh"))
    of "reconcile":
      ensureNotRoot()
      doReconcile(dotDir, homeDir, options, getDotFiles("dotty.sh"))
    of "rootStatus":
      ensureRoot()
      ensureRootFileOwnership(dotDir)
      doStatus(dotDir, homeDir, options, getDotFiles("dottyRoot.sh"))
    of "rootReconcile":
      ensureRoot()
      ensureRootFileOwnership(dotDir)
      doReconcile(dotDir, homeDir, options, getDotFiles("dottyRoot.sh"))
    else:
      die "Subcommand not recognized"
