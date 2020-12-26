import os
import strutils
proc xdgCfg*: string =
  if getEnv("XDG_CONFIG_HOME") != "":
    return getEnv("XDG_CONFIG_HOME")
  else:
    return joinPath(getEnv("HOME"), ".config")

proc xdgData*: string =
  if getEnv("XDG_DATA_HOME") != "":
    return getEnv("XDG_DATA_HOME")
  else:
    return joinPath(getEnv("HOME"), ".local/share")

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
