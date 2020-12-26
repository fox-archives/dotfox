import os
import system
import strutils
import parsetoml

import "./util"

proc doStatus(dotDir: string, homeDir: string, dotFiles: seq[string]) =
  for i, file in dotFiles:
    if symlinkExists(file):
      if symlinkResolvedProperly(dotDir, homeDir, file):
        # symlinks pointing to a file or folder may or may not have a trailing slash. skip ones that do for consistency
        if endsWith(expandSymlink(file), '/'):
            echo "[OK_SLASH] " & file
        else:
            echo "[OK]       " & file
      else:
        echo "[BROKEN_S] " & file
        echo "          -> " & expandSymlink(file)
    elif fileExists(file):
        echo "[ROGUE_F]  " & file
    elif dirExists(file):
        echo "[ROGUE_D]  " & file
    # nothing exists, we create it
    else:
        let rel = file[len(homeDir) .. ^1]
        if fileExists(joinPath(dotDir, rel)):
          echo "[TOLINK_F] " & file
        elif dirExists(joinPath(dotDir, rel)):
          echo "[TOLINK_D] " & file
        else:
          echo "[MISSING]  " & file

proc doReconcile(dotDir: string, homeDir: string, dotFiles: seq[string]) =
  for i, file in dotFiles:
    # ensure directory exists
    createDir(parentDir(file))

    if symlinkExists(file):
        if symlinkResolvedProperly(dotDir, homeDir, file):
          # everything OK except the trailing slash; fix it
          if endsWith(expandSymlink(file), '/'):
              echo  "FIX OK_SLASH: " & file
              let temp = expandSymlink(file)
              removeFile(file)
              createSymlink(rts(temp), file)
        else:
          echo "FIX BROKEN_S: " & file
          removeFile(file)
          createSymlink(getRealDot(dotDir, homeDir, file), file)
    elif fileExists(file):
        let real = getRealDot(dotDir, homeDir, file)

        if fileExists(real):
          let fileContents = readFile(file)
          let realContents = readFile(real)

          if fileContents == realContents:
            removeFile(file)
            createSymlink(real, file)
          else:
            echo "SKIP ROGUE_F_F Path conflict: Remove the outdated and try again"
            echo "             -> " & file & " (file)"
            echo "             -> " & real & " (file)"
        elif dirExists(real):
          echo "SKIP ROGUE_F_D Path conflict: Remove the outdated and try again"
          echo "             -> " & file & " (file)"
          echo "             -> " & real & " (directory)"
        else:
          echo "FIX ROGUE_F_M  " & file
          # ensure directory
          createDir(parentDir(real))

          # file doesn't exist on other side. move it
          moveFile(file, real)
          createSymlink(real, file)

    elif dirExists(file):
        let real = getRealDot(dotDir, homeDir, file)

        if fileExists(real):
          echo "SKIP ROGUE_D_F Path conflict: Remove the outdated and try again"
          echo "             -> " & file & " (directory)"
          echo "             -> " & real & " (file)"
        elif dirExists(real):
          if dirLength(file) == 0:
            echo "FIX ROGUE_D  " & file
            removeDir(file)
            createSymlink(joinPath(dotDir, getRel(homeDir, file)), file)
          # TODO: do some merging or whatever
          else:
            echo "SKIP ROGUE_D_D Path conflict: Remove the outdated and try again"
            echo "             -> " & file & " (directory)"
            echo "             -> " & real & " (directory)"
        else:
          echo "FIX ROGUE_D_M  " & file
          # ensure directory
          createDir(parentDir(real))

          # file doesn't exist on other side. move it
          try:
            copyDirWithPermissions(file, real)
            removeDir(file)
            createSymlink(real, file)
          except Exception:
            echo "Error: ROGUE_D_M Could not copy folder"

    else:
        createSymlink(joinPath(dotDir, getRel(homeDir, file)), file)

let toml = parsetoml.parseFile(joinPath(getConfigDir(), "dotty/dotty.toml"))
let cfg = toml["config"]

if paramCount() < 1:
  echo "Error: Expected subcommand. Exiting"
  quit 1

let dotDir = expandTilde(cfg["dotDir"].getStr())
let homeDir = expandTilde(cfg["homeDir"].getStr())

case paramStr(1):
  of "status":
    doStatus(dotDir, homeDir,  getDotFiles())
  of "reconcile":
    doReconcile(dotDir, homeDir, getDotFiles())
  else:
    echo "Error: Subcommand not found. Exiting"
    quit 1
