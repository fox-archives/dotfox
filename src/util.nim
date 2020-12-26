import os
import osproc
import sequtils
import strutils

proc getDotFiles*(): seq[string] =
  let cfg = joinPath(getConfigDir(), "dotty", "dotty.nim")
  if not fileExists(cfg):
    echo "dotty.nim not found at '" & cfg & "'. Create one"
    quit 1

  let output = execProcess("nim", args=["--hints:off", "r", cfg], options={poUsePath})
  var dotFiles = newSeq[string]()
  for str in filter(output.split('\n'), proc(str: string): bool = not isEmptyOrWhitespace(str)):
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

# for when dotFile / dotFolder doesn't exist, but the symlink points there
proc createRel*(dotFile: string) =
  echo "do thing"

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
