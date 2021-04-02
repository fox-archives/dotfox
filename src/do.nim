import os
import system
import strutils
import "./util"
import strformat
import posix

# for each higher order function (ex. runSymlinkDir), the first word (e.g. Symlink) represents the type of file
# located in the home / destination folder. The Second word (ex. Dir) represents the type of
# file that exists in the dotfile repo
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

    let real = getRealDot(dotDir, homeDir, file)

    # if running as root, ensure all 'real' files
    # in dotDir are owned by root. prevent security issues
    # since non-root users would be able to write to a file
    # which could eventually be ran in a root context
    if geteuid() == 0:
      try:
        if not hasRootOwnership(real):
          setRootOwnership(real)
      except Exception:
        logError getCurrentExceptionMsg()
        echoStatus("SKIP", file)
        continue

    if symlinkExists(file):
      if fileExists(real):
        runSymlinkFile(file, real)
      elif dirExists(real):
        runSymlinkDir(file, real)
      else:
        runSymlinkNull(file, real)

    elif fileExists(file):
      if fileExists(real):
        runFileFile(file, real)
      elif dirExists(real):
        runFileDir(file, real)
      else:
        runFileNull(file, real)

    elif dirExists(file):
      if fileExists(real):
        runDirFile(file, real)
      elif dirExists(real):
        runDirDir(file, real)
      else:
        runDirNull(file, real)

    else:
      if fileExists(real):
        runNullFile(file, real)
      elif dirExists(real):
        runNullDir(file, real)
      else:
        runNullNull(file, real)


proc doStatus*(dotDir: string, homeDir: string, dotFiles: seq[string]) =
  proc runSymlinkFile(file: string, real: string): void =
    if symlinkResolvedProperly(dotDir, homeDir, file):
      if endsWith(expandSymlink(file), '/'):
        echoStatus("OK_S", file)
      else:
        echoStatus("OK", file)
    else:
      echoStatus("Y_SYM_FILE", file)

  proc runSymlinkDir(file: string, real: string): void =
    if symlinkResolvedProperly(dotDir, homeDir, file):
      if endsWith(expandSymlink(file), '/'):
        echoStatus("OK_S", file)
      else:
        echoStatus("OK", file)
    else:
      echoStatus("Y_SYM_DIR", file)

  proc runSymlinkNull(file: string, real: string): void =
    echoStatus("M_SYM_NULL", file)

  proc runFileFile(file: string, real: string): void =
    echoStatus("E_FILE_FILE", file)

  proc runFileDir(file: string, real: string): void =
    echoStatus("E_FILE_DIR", file)

  proc runFileNull(file: string, real: string): void =
    echoStatus("Y_FILE_NULL", file)

  proc runDirFile(file: string, real: string): void =
    echoStatus("E_DIR_FILE", file)

  proc runDirDir(file: string, real: string): void =
    echoStatus("E_DIR_DIR", file)

  proc runDirNull(file: string, real: string): void =
    echoStatus("Y_DIR_NULL", file)

  proc runNullFile(file: string, real: string): void =
    echoStatus("Y_NULL_FILE", file)

  proc runNullDir(file: string, real: string): void =
    echoStatus("Y_NULL_DIR", file)

  proc runNullNull(file: string, real: string): void =
    echoStatus("M_NULL_NULL", file)

  doAbstract(
    dotDir,
    homeDir,
    dotFiles,
    runSymlinkFile,
    runSymlinkDir,
    runSymlinkNull,
    runFileFile,
    runFileDir,
    runFileNull,
    runDirFile,
    runDirDir,
    runDirNull,
    runNullFile,
    runNullDir,
    runNullNull
  )

# proc doRootStatus*(dotDir: string, homeDir: string, dotFiles: seq[string]) =
#   # mostly similar to doRootStatus except for symlinkAny

#   proc runSymlinkFile(file: string, real: string): void =
#     echoStatus("Y_SYM_FILE", file)

#   proc runSymlinkDir(file: string, real: string): void =
#     echoStatus("Y_SYM_DIR", file)

#   proc runSymlinkNull(file: string, real: string): void =
#     echoStatus("Y_SYM_NULL", file)

#   proc runFileFile(file: string, real: string): void =
#     echoStatus("E_FILE_FILE", file)

#   proc runFileDir(file: string, real: string): void =
#     echoStatus("E_FILE_DIR", file)

#   proc runFileNull(file: string, real: string): void =
#     echoStatus("Y_FILE_NULL", file)

#   proc runDirFile(file: string, real: string): void =
#     echoStatus("E_DIR_FILE", file)

#   proc runDirDir(file: string, real: string): void =
#     echoStatus("E_DIR_DIR", file)

#   proc runDirNull(file: string, real: string): void =
#     echoStatus("Y_DIR_NULL", file)

#   proc runNullFile(file: string, real: string): void =
#     echoStatus("Y_NULL_FILE", file)

#   proc runNullDir(file: string, real: string): void =
#     echoStatus("Y_NULL_DIR", file)

#   proc runNullNull(file: string, real: string): void =
#     echoStatus("M_NULL_NULL", file)

#   doAbstract(
#     dotDir,
#     homeDir,
#     dotFiles,
#     runSymlinkFile,
#     runSymlinkDir,
#     runSymlinkNull,
#     runFileFile,
#     runFileDir,
#     runFileNull,
#     runDirFile,
#     runDirDir,
#     runDirNull,
#     runNullFile,
#     runNullDir,
#     runNullNull
#   )

proc doReconcile*(dotDir: string, homeDir: string, dotFiles: seq[string]) =
  proc runSymlinkFile(file: string, real: string) =
    if symlinkResolvedProperly(dotDir, homeDir, file):
      # if destination has an extraneous forward slash,
      # automatically remove it
      if endsWith(expandSymlink(file), '/'):
        let temp = expandSymlink(file)
        removeFile(file)
        createSymlink(rts(temp), file)
    else:
      removeFile(file)
      createSymlink(getRealDot(dotDir, homeDir, file), file)

  proc runSymlinkDir(file: string, real: string) =
    if symlinkResolvedProperly(dotDir, homeDir, file):
      # if destination has a spurious slash, automatically
      # remove it
      if endsWith(expandSymlink(file), '/'):
        let temp = expandSymlink(file)
        removeFile(file)
        createSymlink(rts(temp), file)
    else:
      removeFile(file)
      createSymlink(getRealDot(dotDir, homeDir, file), file)

  proc runSymlinkNull(file: string, real: string) =
    echoStatus("M_SYM_NULL", file)
    echoPoint(fmt"{file} (symlink)")
    echoPoint(fmt"{real} (nothing here)")

  proc runFileFile(file: string, real: string) =
    let fileContents = readFile(file)
    let realContents = readFile(real)

    if fileContents == realContents:
      removeFile(file)
      createSymlink(real, file)
    else:
      echoStatus("E_FILE_FILE", file)
      echoPoint(fmt"{file} (file)")
      echoPoint(fmt"{real} (file)")

  proc runFileDir(file: string, real: string) =
    echoStatus("E_FILE_DIR", file)
    echoPoint(fmt"{file} (file)")
    echoPoint(fmt"{real} (directory)")

  proc runFileNull (file: string, real: string) =
    echoStatus("E_FILE_NULL", file)
    echoPoint("Automatically fixed")

    createDir(parentDir(real))

    # file doesn't exist on other side. move it
    moveFile(file, real)
    createSymlink(real, file)

  proc runDirFile (file: string, real: string) =
    echoStatus("E_DIR_FILE", file)
    echoPoint(fmt"{file} (directory)")
    echoPoint(fmt"{real} (file)")

  proc runDirDir (file: string, real: string) =
    if dirLength(file) == 0:
      echoStatus("E_DIR_DIR", file)
      echoPoint("Automatically fixed")

      removeDir(file)
      createSymlink(joinPath(dotDir, getRel(homeDir, file)), file)
    # TODO
    # elif dirLength(real) == 0:
    else:
      echoStatus("E_DIR_DIR", file)
      echoPoint(fmt"{file} (directory)")
      echoPoint(fmt"{file} (directory)")

  proc runDirNull (file: string, real: string) =
    # ensure directory
    createDir(parentDir(real))

    # file doesn't exist on other side. move it
    try:
      copyDirWithPermissions(file, real)
      removeDir(file)
      createSymlink(real, file)

      echoStatus("E_DIR_NULL", file)
      echoPoint("Automatically fixed")
    except Exception:
      logError getCurrentExceptionMsg()
      echoStatus("E_DIR_NULL", file)
      echoPoint("Error: Could not copy folder")

  proc runNullAny(file: string, real: string) =
    createSymlink(joinPath(dotDir, getRel(homeDir, file)), file)

  doAbstract(
    dotDir,
    homeDir,
    dotFiles,
    runSymlinkFile,
    runSymlinkDir,
    runSymlinkNull,
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

# proc doRootReconcile*(dotDir: string, homeDir: string, dotFiles: seq[string]) =
#   doReconcile(dotDir, homeDir, dotFiles)
