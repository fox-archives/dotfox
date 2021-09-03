import os
import system
import strutils
import strformat
import "./util"

# For each higher order function (ex. runSymlinkDir), the first word (e.g. Symlink) represents the type of file
# located in the home / destination folder. The Second word (ex. Dir) represents the type of
# file that exists in the dotfile repo
proc doAbstract(
  dotDir: string,
  homeDir: string,
  options: Options,
  dotfiles: seq[array[2, string]],
  runSymlinkSymlink: proc (dotfile: string, srcFile: string, options: Options),
  runSymlinkFile: proc (dotfile: string, srcFile: string, options: Options),
  runSymlinkDir: proc (dotfile: string, srcFile: string, options: Options),
  runSymlinkNull: proc (dotfile: string, srcFile: string),
  runFileFile: proc(dotfile: string, srcFile: string),
  runFileDir: proc(dotfile: string, srcFile: string),
  runFileNull: proc (dotfile: string, srcFile: string),
  runDirFile: proc(dotfile: string, srcFile: string),
  runDirDir: proc(dotfile: string, srcFile: string),
  runDirNull: proc(dotfile: string, srcFile: string),
  runNullFile: proc(dotfile: string, srcFile: string),
  runNullDir: proc(dotfile: string, srcFile: string),
  runNullNull: proc(dotfile: string, srcFile: string)
) =
  for i, files in dotfiles:
    let srcFile = files[0]
    let destFile = files[1]

    try:
      createDir(parentDir(destFile))

      if symlinkExists(destFile):
        if symlinkExists(srcFile):
          runSymlinkSymlink(destFile, srcFile, options)
        elif fileExists(srcFile):
          runSymlinkFile(destFile, srcFile, options)
        elif dirExists(srcFile):
          runSymlinkDir(destFile, srcFile, options)
        else:
          runSymlinkNull(destFile, srcFile)

      elif fileExists(destFile):
        if fileExists(srcFile):
          runFileFile(destFile, srcFile)
        elif dirExists(srcFile):
          runFileDir(destFile, srcFile)
        else:
          runFileNull(destFile, srcFile)

      elif dirExists(destFile):
        if fileExists(srcFile):
          runDirFile(destFile, srcFile)
        elif dirExists(srcFile):
          runDirDir(destFile, srcFile)
        else:
          runDirNull(destFile, srcFile)

      else:
        if fileExists(srcFile):
          runNullFile(destFile, srcFile)
        elif dirExists(srcFile):
          runNullDir(destFile, srcFile)
        else:
          runNullNull(destFile, srcFile)
    except Exception:
      logError &"Unhandled exception raised\n{getCurrentExceptionMsg()}"
      printStatus("SKIP", destFile)
  echo "Done."

proc doStatus*(dotDir: string, homeDir: string, options: Options, dotfiles: seq[array[2, string]]) =
  proc runSymlinkSymlink(file: string, srcFile: string, options: Options): void =
    # This is possible if dotty does it's thing correctly, but
    # the user replaces the file/directory in dotDir with a symlink
    # to something else. It is an error, even if the symlink resolves
    # properly, and it should not be possible in normal circumstances
    printStatus("ERR_SYM_SYM", file)
    printHint("(not fixable)")

  proc runSymlinkFile(file: string, srcFile: string, options: Options): void =
    if symlinkResolvedProperly(file, srcFile):
      if endsWith(expandSymlink(file), '/'):
        if options.showOk:
          printStatus("OK/", file)
      else:
        if options.showOk:
          printStatus("OK", file)
    else:
      printStatus("ERR_SYM_FILE", file)
      # Possibly fixable, see reasoning in runSymlinkDir()
      printHint("(possibly fixable)")

  proc runSymlinkDir(file: string, srcFile: string, options: Options): void =
    if symlinkResolvedProperly(file, srcFile):
      if endsWith(expandSymlink(file), '/'):
        if options.showOk:
          printStatus("OK/", file)
      else:
        if options.showOk:
          printStatus("OK", file)
    else:
      printStatus("ERR_SYM_DIR", file)

  proc runSymlinkNull(file: string, srcFile: string): void =
    printStatus("ERR_SYM_NULL", file)
    printHint(fmt"{file} (symlink)")
    printHint(fmt"{srcFile} (nothing here)")
    printHint(fmt"Did you forget to create your actual dotfile at '{srcFile}'?")
    printHint("(not fixable)")

  proc runFileFile(file: string, srcFile: string): void =
    printStatus("ERR_FILE_FILE", file)
    printHint(fmt"{file} (file)")
    printHint(fmt"{srcFile} (file)")
    printHint("(possibly fixable)")

  proc runFileDir(file: string, srcFile: string): void =
    printStatus("ERR_FILE_DIR", file)
    printHint(fmt"{file} (file)")
    printHint(fmt"{srcFile} (directory)")
    printHint("(not fixable)")

  proc runFileNull(file: string, srcFile: string): void =
    printStatus("ERR_FILE_NULL", file)
    printHint("(fixable)")

  proc runDirFile(file: string, srcFile: string): void =
    printStatus("ERR_DIR_FILE", file)
    printHint(fmt"{file} (directory)")
    printHint(fmt"{srcFile} (file)")
    printHint("(not fixable)")

  proc runDirDir(file: string, srcFile: string): void =
    printStatus("ERR_DIR_DIR", file)
    printHint(fmt"{file} (directory)")
    printHint(fmt"{srcFile} (directory)")
    printHint("Remove the directory that has the older contents")
    printHint("(possibly fixable)")

  proc runDirNull(file: string, srcFile: string): void =
    printStatus("ERR_DIR_NULL", file)
    printHint("(fixable)")

  proc runNullFile(file: string, srcFile: string): void =
    printStatus("ERR_NULL_FILE", file)
    printHint("(fixable)")

  proc runNullDir(file: string, srcFile: string): void =
    printStatus("ERR_NULL_DIR", file)
    printHint("(fixable)")

  proc runNullNull(file: string, srcFile: string): void =
    printStatus("ERR_NULL_NULL", file)
    printHint(fmt"Did you forget to create your actual dotfile at '{srcFile}'?")
    printHint("(not fixable)")

  doAbstract(
    dotDir,
    homeDir,
    options,
    dotfiles,
    runSymlinkSymlink,
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

proc doReconcile*(dotDir: string, homeDir: string, options: Options,
    dotfiles: seq[array[2, string]]) =
  proc runSymlinkSymlink(file: string, srcFile: string, options: Options): void =
    discard # not fixable

  proc runSymlinkFile(file: string, srcFile: string, options: Options) =
    if symlinkResolvedProperly(file, srcFile):
      # If the destination has an extraneous forward slash,
      # automatically remove it
      if endsWith(expandSymlink(file), '/'):
        let temp = expandSymlink(file)
        removeFile(file)
        createSymlink(rts(temp), file)

  proc runSymlinkDir(file: string, srcFile: string, options: Options) =
    if symlinkResolvedProperly(file, srcFile):
      # If the destination has a spurious slash, automatically remove it
      if endsWith(expandSymlink(file), '/'):
        let temp = expandSymlink(file)
        removeFile(file)
        createSymlink(rts(temp), file)

  proc runSymlinkNull(file: string, srcFile: string) =
    discard # not fixable

  proc runFileFile(file: string, srcFile: string) =
    let fileContents = readFile(file)
    let srcFileContents = readFile(srcFile)

    if fileContents == srcFileContents:
      removeFile(file)
      createSymlink(srcFile, file)
    else:
      discard # not fixable

  proc runFileDir(file: string, srcFile: string) =
    discard # not fixable

  proc runFileNull (file: string, srcFile: string) =
    printStatus("ERR_FILE_NULL", file)
    # TODO: make auto fix more consistent
    printHint("Automatically fixed")

    createDir(parentDir(srcFile))

    # The file doesn't exist on other side. Move it
    moveFile(file, srcFile)
    createSymlink(srcFile, file)

  proc runDirFile (file: string, srcFile: string) =
    discard # not fixable

  # Swapped
  proc runDirNull (file: string, srcFile: string) =
    # Ensure directory hierarchy exists
    createDir(parentDir(srcFile))

    # The file doesn't exist on other side. Move it
    try:
      printStatus("ERR_DIR_NULL", file)
      printHint("Automatically fixed")

      copyDirWithPermissions(file, srcFile)
      removeDir(file)
      createSymlink(srcFile, file)
    except Exception:
      logError getCurrentExceptionMsg()
      printStatus("ERR_DIR_NULL", file)
      # TODO: elaborate
      printHint("Error: Could not copy folder")

  # Swapped
  proc runDirDir (file: string, srcFile: string) =
    if dirLength(file) == 0:
      printStatus("ERR_DIR_DIR", file)
      printHint("Automatically fixed")

      removeDir(file)
      createSymlink(srcFile, file)
    elif dirLength(srcFile) == 0:
      printStatus("ERR_DIR_DIR", file)
      printHint("Automatically fixed")

      removeDir(srcFile)
      runDirNull(file, srcFile)
    else:
      discard # not fixable

  proc runNullFile(file: string, srcFile: string) =
    createSymlink(srcFile, file)

  proc runNullDir(file: string, srcFile: string) =
    createSymlink(srcFile, file)

  proc runNullNull(file: string, srcFile: string) =
    discard # not fixable

  doAbstract(
    dotDir,
    homeDir,
    options,
    dotfiles,
    runSymlinkSymlink,
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

proc doDebug*(dotDir: string, homeDir: string, options: Options,
    dotfiles: seq[array[2, string]]) =
  for file in dotfiles:
    echo file
