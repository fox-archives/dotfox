import os
import system
import strutils
import "./util"

proc doAbstract(
  dotDir: string,
  homeDir: string,
  dotFiles: seq[string],
  runSymlinkFile: proc (dotFile: string, real: string),
  runSymlinkDir: proc (dotFile: string, real: string),
  runSymlinkNull: proc (dotFile: string, real: string),
  runFileFile: proc(dotFile: string, real: string),
  runFileDir: proc(dotFile: string, real: string),
  runFileNull: proc (dotFile: string, real: string),
  runDirFile: proc(dotFile: string, real: string),
  runDirDir: proc(dotFile: string, real: string),
  runDirNull: proc(dotFile: string, real: string),
  runNullFile: proc(dotFile: string, real: string),
  runNullDir: proc(dotFile: string, real: string),
  runNullNull: proc(dotFile: string, real: string)
) =
  for i, file in dotFiles:
    createDir(parentDir(file))

    if symlinkExists(file):
      let real = getRealDot(dotDir, homeDir, file)
      if fileExists(real):
        runSymlinkFile(file, real)
      elif dirExists(real):
        runSymlinkDir(file, real)
      else:
        runSymlinkNull(file, real)

    elif fileExists(file):
        let real = getRealDot(dotDir, homeDir, file)
        if fileExists(real):
          runFileFile(file, real)
        elif dirExists(real):
          runFileDir(file, real)
        else:
          runFileNull(file, real)

    elif dirExists(file):
        let real = getRealDot(dotDir, homeDir, file)
        if fileExists(real):
          runDirFile(file, real)
        elif dirExists(real):
          runDirDir(file, real)
        else:
          runDirNull(file, real)

    else:
      let real = getRealDot(dotDir, homeDir, file)
      if fileExists(real):
        runNullFile(file, real)
      elif dirExists(real):
        runNullDir(file, real)
      else:
        runNullNull(file, real)

proc doStatus*(dotDir: string, homeDir: string, dotFiles: seq[string]) =
  proc runSymlinkAny(file: string, real: string) =
    if symlinkResolvedProperly(dotDir, homeDir, file):
      if endsWith(expandSymlink(file), '/'):
        echo "[VALID_SLASH] " & file
      # symlink (valid), ∅
      elif not fileExists(getRealDot(dotDir, homeDir, file)) and not dirExists(getRealDot(dotDir, homeDir, file)):
        echo "[VALID_MISS] " & file
      # symlink (valid), dir
      # symlink (valid), file
      else:
        echo "[OK]         " & file
    # symlink (invalid), -
    # symlink (invalid), ∅
    # symlink (invalid), dir
    # symlink (invalid), file
    else:
      echo "[INVALID] " & file
      echo "          -> " & real

  proc runFileAny(file: string, real: string) =
    echo "[ROGUE_F]  " & file

  proc runDirAny(file: string, real: string) =
    echo "[ROGUE_D]  " & file

  proc runNullFile(file: string, real: string) =
    echo "[TOLINK_F] " & file

  proc runNullDir(file: string, real: string) =
    echo "[TOLINK_D] " & file

  proc runNullNull(file: string, real: string) =
    echo "[MISSING]  " & file

  doAbstract(
    dotDir,
    homeDir,
    dotFiles,
    runSymlinkAny,
    runSymlinkAny,
    runSymlinkAny,
    runFileAny,
    runFileAny,
    runFileAny,
    runDirAny,
    runDirAny,
    runDirAny,
    runNullFile,
    runNullDir,
    runNullNull
  )


proc doReconcile*(dotDir: string, homeDir: string, dotFiles: seq[string]) =
  proc runSymlinkAny(file: string, real: string) =
    if symlinkResolvedProperly(dotDir, homeDir, file):
      # transform trailing slash for consistency
      if endsWith(expandSymlink(file), '/'):
          echo  "FIX OK_SLASH: " & file
          let temp = expandSymlink(file)
          removeFile(file)
          createSymlink(rts(temp), file)
    else:
      echo "FIX BROKEN_S: " & file
      removeFile(file)
      createSymlink(getRealDot(dotDir, homeDir, file), file)

  proc runFileFile(file: string, real: string) =
    let fileContents = readFile(file)
    let realContents = readFile(real)

    if fileContents == realContents:
      removeFile(file)
      createSymlink(real, file)
    else:
      echo "SKIP ROGUE_F_F Path conflict: Remove the outdated and try again"
      echo "             -> " & file & " (file)"
      echo "             -> " & real & " (file)"
  proc runFileDir(file: string, real: string) =
    echo "SKIP ROGUE_F_D Path conflict: Remove the outdated and try again"
    echo "             -> " & file & " (file)"
    echo "             -> " & real & " (directory)"
  proc runFileNull (file: string, real: string) =
    echo "FIX ROGUE_F_M  " & file
    # ensure directory
    createDir(parentDir(real))

    # file doesn't exist on other side. move it
    moveFile(file, real)
    createSymlink(real, file)

  proc runDirFile (file: string, real: string) =
    echo "SKIP ROGUE_D_F Path conflict: Remove the outdated and try again"
    echo "             -> " & file & " (directory)"
    echo "             -> " & real & " (file)"
  proc runDirDir (file: string, real: string) =
    if dirLength(file) == 0:
      echo "FIX ROGUE_D  " & file
      removeDir(file)
      createSymlink(joinPath(dotDir, getRel(homeDir, file)), file)
    # TODO: do some merging or whatever
    else:
      echo "SKIP ROGUE_D_D Path conflict: Remove the outdated and try again"
      echo "             -> " & file & " (directory)"
      echo "             -> " & real & " (directory)"
  proc runDirNull (file: string, real: string) =
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

  proc runNullAny(file: string, real: string) =
    createSymlink(joinPath(dotDir, getRel(homeDir, file)), file)

  doAbstract(
    dotDir,
    homeDir,
    dotFiles,
    runSymlinkAny,
    runSymlinkAny,
    runSymlinkAny,
    runFileFile,
    runFileDir,
    runFileNull,
    runDirFile,
    runDirDir,
    runDirNull,
    runNullAny,
    runNullAny,
    runNullAny
  )
