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
    action*: string
    configFile*: string
    deployment*: string

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

proc printStatus*(status: string, file: string): void =
  # Print the status code for the particular file
  let s = fmt"[{status}]"
  echo fmt"{s:<16}" & file

proc printHint*(str: string): void =
  # Print a hint for a particular file, but indented so the output is more clear
  echo fmt"                -> {str}"  

proc hasAllRootFiles*(parentDir: string): bool =
  # Determine if all child files and subdirectories of the particular directory are owned by root.
  # We use this as a security precaution when operating on `sudo` mode. This is because we don't want
  # to copy files that are user-writable to a root level location
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
  if not fileOwnedByRoot(parentDir):
    echo parentDir
    allRootFiles = false

  walk(parentDir)
  return allRootFiles

proc getDotfileList*(deploymentStr: string): seq[string] =
  # Execute deployment, returning all their standard output, concatenated
  let oldCurrentDir = getCurrentDir()
  setCurrentDir(joinPath(getConfigDir(), "dotty", "deployments"))

  var deployment = ""

  if not isAbsolute(deploymentStr):
    deployment = joinPath(getConfigDir(), "dotty", "deployments", deploymentStr)
  else:
    deployment = deploymentStr

  if not fileExists(deployment):
    die fmt"Deployment file '{deployment}' not found"

  let cmdResult = execCmdEx(deployment)
  if cmdResult.exitCode != 0:
    stdout.write cmdResult.output
    die fmt"Executing {deployment} failed"

  var dotfiles = newSeq[string]()
  for str in filter(cmdResult.output.split('\n'), proc(
      str: string): bool = not isEmptyOrWhitespace(str)):
    dotfiles.add(str)

  setCurrentDir(oldCurrentDir)
  return dotfiles


proc rts*(str: string): string =
  # Remove trailing slash
  if endsWith(str, '/'):
    return str[0 .. ^2]

  return str

proc getRel*(homeDir: string, dotfile: string): string =
  # From dotfile (in homeDir), get relative path
  return dotfile[len(homeDir) .. ^1]

proc getRealDot*(dotDir: string, homeDir: string, dotfile: string): string =
  # From dotfile (in homeDir), get the real path that's in dotDir
  return joinPath(dotDir, getRel(homeDir, dotfile))

proc symlinkCreatedByDotty*(dotDir: string, homeDir: string,
    symlinkFile: string): bool =
  # Determine if the symlink has been created by us. We it is if it points to somewhere in `dotDir`
  if startsWith(symlinkFile, dotDir):
    return true
  return false

proc symlinkResolvedProperly*(dotDir: string, homeDir: string,
    dotfile: string): bool =
  # Test if the symlink in homeDir actually points to corresponding one in dotfile. This
  # assumes the symlink exists
  if rts(expandSymlink(dotfile)) == joinPath(dotDir, getRel(homeDir, dotfile)):
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
    associated with it when using the 'status' subcommand
  --roots
    Manage the dotfiles for the root user
  --deployment
    Specify specific deployment to read dotfiles from. This defaults to 'dotty.sh'
    in your config directory
  --config
    Set location of config file. This defaults to 'config.toml'
    in your config directory

Usage:
  dotty --show-ok=false status
  sudo dotty reconcile --root
"""