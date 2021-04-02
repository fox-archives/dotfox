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
    die "Must not be running as root"

proc hasRootOwnership*(path: string): bool =
  var info: Stat
  let err = stat(path, info)
  if err != 0:
    raise Exception.newException(fmt"stat error: {path} (does it exist?)")

  let passwd: ptr Passwd = getpwuid(info.st_uid)
  return passwd[].pw_uid == 0

proc setRootOwnership*(path: string): void =
  var info: Stat
  let err = stat(path, info)
  if err != 0:
    raise Exception.newException(fmt"stat error: {path} (does it exist?)")

  # keep same gid since some systems like OpenSUSE have
  # a different group ownership
  let success = chown(path.cstring, 0.Uid, info.st_gid)
  if success != 0:
    raise Exception.newException(fmt"chown error: {path}")

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
