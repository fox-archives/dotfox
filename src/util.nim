import os
import osproc
import sequtils
import strutils
import std/terminal
import strformat
import posix

type
  Options* = object
    showOk*: bool
    isRoot*: bool
    action*: string
    configDir*: string
    deployment*: string

proc logError*(str: string): void {.inline.} =
  writeLine stderr, fmt"{ansiForegroundColorCode(fgRed)}Error:{ansiResetCode} {str}"
  flushFile stderr

proc logWarn*(str: string): void {.inline.} =
  writeLine stderr, fmt"{ansiForegroundColorCode(fgYellow)}Info:{ansiResetCode} {str}"
  flushFile stderr

proc logInfo*(str: string): void {.inline.} =
  echo fmt"{ansiForegroundColorCode(fgGreen)}Info:{ansiResetCode} {str}"

proc die*(str: string): void {.noReturn.} =
  logError(fmt"{str}. Exiting")
  quit QuitFailure

proc printStatus*(status: string, file: string): void {.inline.} =
  # Print the status code for the particular file
  let s = fmt"[{status}]"
  echo fmt"{s:<16}" & file

proc printHint*(str: string): void {.inline.} =
  # Print a hint for a particular file, but indented so the output is more clear
  echo fmt"                -> {str}"

proc hasAllRootFiles*(parentDir: string): bool =
  # Determine if all child files and subdirectories of the particular directory are owned by root.
  # We use this as a security precaution when operating on `sudo` mode. This is because we don't want
  # to copy files that are user-writable to a root-level location
  proc fileOwnedByRoot(path: string): bool =
    var info: Stat
    let code = stat(path, info)
    if code != 0:
      logError fmt"Could not stat {path}"
      return false

    # Skip group check
    let fileUsername = getpwuid(info.st_uid).pw_name
    return fileUsername == "root"

  var allRootFiles = true

  # Check root
  if not fileOwnedByRoot(parentDir):
    echo parentDir
    allRootFiles = false

  # Check all children of root
  proc walk(path: string) =
    for kind, file in walkDir(path):
      if not fileOwnedByRoot(file):
        echo file
        allRootFiles = false

      if kind == PathComponent.pcDir:
        walk(file)

  walk(parentDir)
  return allRootFiles

proc getDotfileList*(options: Options): seq[array[2, string]] =
  # Execute deployment, returning all their standard output, concatenated
  let oldCurrentDir = getCurrentDir()
  setCurrentDir(joinPath(options.configDir, "deployments"))

  var deployment = ""

  if not isAbsolute(options.deployment):
    deployment = joinPath(options.configDir, "deployments", options.deployment)
  else:
    deployment = options.deployment

  if not fileExists(deployment):
    die fmt"Deployment file '{deployment}' not found"

  let cmdResult = execCmdEx(deployment)
  if cmdResult.exitCode != 0:
    stdout.write cmdResult.output
    die fmt"Executing {deployment} failed"

  var dotfiles = newSeq[array[2, string]]()
  for line in filter(cmdResult.output.split('\n'), proc(line: string): bool = not isEmptyOrWhitespace(line)):
    if line[0] == '#':
      continue

    let lineParts = line.split(':')
    if len(lineParts) < 1:
      die fmt"Line '{line}' must have one colon, but none were found"

    let prefix = lineParts[0]
    if prefix == "symlink":
      if len(lineParts) != 3:
        die fmt"Symlink prefix on line '{line}' must have 3 elements"

      dotfiles.add([lineParts[1], lineParts[2]])
    else:
        die fmt"Prefix  '{prefix}' in line '{line}' not supported"


  setCurrentDir(oldCurrentDir)
  return dotfiles

proc rts*(str: string): string {.inline.} =
  # Remove trailing slash
  if endsWith(str, '/'):
    return str[0 .. ^2]

  return str

proc symlinkResolvedProperly*(destFile: string, srcFile: string): bool =
  # Test if the symlink in homeDir actually points to corresponding one in dotfile. This
  # assumes the symlink exists
  if rts(expandSymlink(destFile)) == rts(srcFile):
    return true
  else:
    return false

proc dirLength*(dir: string): int =
  # Determine how many files or folders are contained are in a directory
  var len = 0
  for kind, path in walkDir(dir):
    len = len + 1

  return len

proc parseBoolFlag*(flag: string): bool =
  if(flag == "true"):
    return true
  elif(flag == "false"):
    return false
  else:
    die fmt"Value '{flag}' not understood. Use 'true' or 'false'"

proc writeHelp*() =
  echo """dotfox

Usage: dotfox [flags] [subcommand]

Subcommands:
  status
    Views the status of user dotfiles
  deploy
    Symlinks dotfiles to proper location and attempts to autofix mismatches

Flags:
  --help, -h
  --version, -v
  --show-ok
    Only prints information associated with a file if there is an error
    associated with it when using the 'status' subcommand
  --root
    Manage the dotfiles for the root user
  --deployment
    Specify specific deployment to read dotfiles from. This defaults to 'dotfox.sh'
    in your config directory
  --config-dir
    Set location of config directory. This defaults to '~/.config/dotfox'

Usage:
  dotfox --show-ok=false status
  sudo dotfox deploy --root
"""
