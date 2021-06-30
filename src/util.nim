import os
import osproc
import sequtils
import strutils
import terminal
import strformat
import posix

type
  Options* = object
    showOk*: bool
    isRoot*: bool
    interactive*: bool
    action*: string
    configFile*: string
    files*: seq[string]  

proc writeHelp*() =
  echo """Dotty

Usage: dotty [flags] [subcommand]

Subcommands:
  status
    Views the status of user dotfiles
  reconcile
    Symlinks dotfiles to proper location and attempts to autofix mismatches

Flags:
  --help, -h
  --version, -v
  --show-ok
    Only prints information associated with a file if there is an error
    associated with it
  --root
    Manage the dotfiles for the root user
  --interactive, -i
    Enable interactive mode, which shows a prompt before each action (linking, deletion)

Usage:
  dotty --show-ok=false status
  sudo dotty reconcile --root
"""

proc writeVersion*() =
  # TODO
  echo "0.5.0"

proc logError*(str: string): void =
  echo fmt"{ansiForegroundColorCode(fgRed)}Error: {str}"
  resetAttributes()

proc logWarn*(str: string): void =
  echo fmt"{ansiForegroundColorCode(fgYellow)}Info: {str}"
  resetAttributes()

proc logInfo*(str: string): void =
  echo fmt"{ansiForegroundColorCode(fgGreen)}Info: {str}"
  resetAttributes()

proc die*(str: string): void {.noReturn.} =
  logError(fmt"{str}. Exiting")
  quit QuitFailure

proc echoStatus*(status: string, file: string): void =
  let s = fmt"[{status}]"
  echo fmt"{s:<16}" & file

proc echoPoint*(str: string): void =
  echo fmt"                -> {str}"  

proc hasAllRootFiles*(dotDir: string): bool =
  proc fileOwnedByRoot(path: string): bool =
    var info: Stat
    let code = stat(path, info)
    if code != 0:
      logError fmt"Could not stat {path}"
      return false

    # skip group check
    let fileUsername = getpwuid(info.st_uid).pw_name
    return fileUsername == "root"

  var allRootFiles = true
  proc walk(path: string) =
    for kind, file in walkDir(path):
      if not fileOwnedByRoot(file):
        echo file
        allRootFiles = false

      if kind == PathComponent.pcDir:
        walk(file)

  # ensure we check root parent
  if not fileOwnedByRoot(dotDir):
    echo dotDir
    allRootFiles = false

  walk(dotDir)
  return allRootFiles

proc getDotFiles*(files: seq[string]): seq[string] =
  var dotFiles = newSeq[string]()

  for file in files:
    let cfg = file
    
    if not fileExists(cfg):
      die fmt"File '{cfg}' not found"

    let cmdResult = execCmdEx(cfg)
    if cmdResult.exitCode != 0:
      stdout.write cmdResult.output
      die fmt"Executing {cfg} failed"

    for str in filter(cmdResult.output.split('\n'), proc(
        str: string): bool = not isEmptyOrWhitespace(str)):
      dotFiles.add(str)

  return dotFiles

# remove trailing slash
proc rts*(str: string): string =
  if endsWith(str, '/'):
    return str[0 .. ^2]
  return str

# from dotfile (in homeDir), get rel path
proc getRel*(homeDir: string, dotFile: string): string =
  let rel = dotFile[len(homeDir) .. ^1]
  return rel

# from dotFile (in homeDir), get real path that's in dotDir
proc getRealDot*(dotDir: string, homeDir: string, dotFile: string): string =
  return joinPath(dotDir, getRel(homeDir, dotFile))

# determine if the symlink is created by us (explanation in getRealDot())
proc symlinkCreatedByDotty*(dotDir: string, homeDir: string,
    symlinkFile: string): bool =
  if startsWith(symlinkFile, dotDir):
    return true
  return false

# test if the symlink in homeDir actually points to corresponding one in dotFile
# assumes the symlink exists
proc symlinkResolvedProperly*(dotDir: string, homeDir: string,
    dotFile: string): bool =
  if rts(expandSymlink(dotFile)) == joinPath(dotDir, getRel(homeDir, dotFile)):
    return true
  else:
    return false

proc dirLength*(dir: string): int =
  var len = 0
  for kind, path in walkDir(dir):
    len = len + 1
  return len

proc parseCategories*(categories: string): seq[string] =
  if categories == "":
    logError "No files was specified. Please specify a category"
    quit QuitFailure
  elif categories.contains(","):
    return categories.split(",")
  else:
    return @[categories]