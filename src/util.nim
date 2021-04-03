import os
import osproc
import sequtils
import strutils
import posix
import terminal
import strformat

proc writeHelp*() =
  echo """Dotty

Usage: dotty [flags] [subcommand]

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

proc writeVersion*() =
  # TODO
  echo "0.2.1"

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
  echo fmt"{s:<14}" & file

proc echoPoint*(str: string): void =
  echo fmt"              -> {str}"

proc ensureRoot*(): void =
  if geteuid() != 0:
    die "Must be running as root"

proc ensureNotRoot*(): void =
  if geteuid() == 0:
    die "Must NOT be running as root"

proc hasAllRootFiles(dotDir: string): bool =
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

proc ensureRootFileOwnership*(dotDir: string): void =
  if not hasAllRootFiles(dotDir):
    echo fmt"Not all files in {dotDir} are owned by root. Fix this"
    quit QuitFailure

proc getDotFiles*(file: string): seq[string] =
  let cfg = joinPath(getConfigDir(), "dotty", file)
  if not fileExists(cfg):
    die fmt"{file} not found at '{cfg}'"

  let result = execCmdEx(cfg)
  if result.exitCode != 0:
    stdout.write result.output
    die fmt"Executing {cfg} failed"

  var dotFiles = newSeq[string]()
  for str in filter(result.output.split('\n'), proc(str: string): bool = not isEmptyOrWhitespace(str)):
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

# test if the symlink in homeDir actually points to corresponding one in dotFile
# assumes the symlink exists
proc symlinkResolvedProperly*(dotDir: string, homeDir: string, dotFile: string): bool =
  if rts(expandSymlink(dotFile)) == getRealDot(dotDir, homeDir, dotFile):
    return true
  else:
    return false

proc dirLength*(dir: string): int =
  var len = 0
  for kind, path in walkDir(dir):
    len = len + 1
  return len
