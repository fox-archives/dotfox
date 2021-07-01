import os
import system
import strutils
import strformat
import "./util"

# for each higher order function (ex. runSymlinkDir), the first word (e.g. Symlink) represents the type of file
# located in the home / destination folder. The Second word (ex. Dir) represents the type of
# file that exists in the dotfile repo
proc doAbstract(
  dotDir: string,
  homeDir: string,
  options: Options,
  dotFiles: seq[string],
  runSymlinkSymlink: proc (dotFile: string, real: string, options: Options),
  runSymlinkFile: proc (dotFile: string, real: string, options: Options),
  runSymlinkDir: proc (dotFile: string, real: string, options: Options),
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
    try:
      createDir(parentDir(file))

      # 'file' and 'dotFile' are synonymous
      # if dotFile is a symlink, it could mean the symlink was created by dotty,
      # or by the user (either could point to a symlink, file, or directory)
      # (ex. ~/.profile -> ~/.config/profile/profile.sh) (created by user)
      # (ex. ~/bin -> ~/.local/bin) (created by user)
      # (ex. ~/.config/profile/profile.sh -> ~/.dots/.config/profile/profile.sh) (created by dotty)

      # we must test if the symlink points to a symlink/file/dir ("real")
      # that has a prefix the same as dotDir. if it does, we dotty created the dotfile. if it
      # does not, the user created the dotfile. to make this check possible (in runSymlinkSymlink,
      # runSymlinkFile, runSymlinkDir, and runSymlinkNull), we HAVE to return
      # "rts(expandSymlink(dotFile))" so we can test the "real" path (rather than defaulting
      # to a ERR_SYM_NULL error with "joinPath(dotDir, getRel(homeDir, dotFile))")
      var real = ""
      if symlinkExists(file):
        # if the symlink expands to a folder, it will append a slash,
        # causing symlinkExists() to fail. rts() rectifies this
        real = rts(expandSymlink(file))
      else:
        real = getRealDot(dotDir, homeDir, file)

      if symlinkExists(file):
        if symlinkExists(real):
          runSymlinkSymlink(file, real, options)
        elif fileExists(real):
          runSymlinkFile(file, real, options)
        elif dirExists(real):
          runSymlinkDir(file, real, options)
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
    except Exception:
      logError &"Unhandled exception raised\n{getCurrentExceptionMsg()}"
      echoStatus("SKIP", file)
  echo "Done."

proc doStatus*(dotDir: string, homeDir: string, options: Options, dotFiles: seq[string]) =
  proc runSymlinkSymlink(file: string, real: string, options: Options): void =
    if symlinkCreatedByDotty(dotDir, homeDir, real):
      # this is possible if dotty does it's thing correctly, but
      # the user replaces the file/directory in dotDir with a symlink
      # to something else. it is an error, even if the symlink resolves
      # properly, and it should not be possible in normal circumstances
      echoStatus("ERR_SYM_SYM", file)
      echoPoint("(not fixable)")
    # symlink created by user
    else:
      # even if symlink does not point to a valid location, we print OK
      # since the symlink is created by the user and we don't track those
      if options.showOk:
        echoStatus("OK_USYM_SYM", file)

  proc runSymlinkFile(file: string, real: string, options: Options): void =
    if symlinkCreatedByDotty(dotDir, homeDir, real):
      if symlinkResolvedProperly(dotDir, homeDir, file):
        if endsWith(expandSymlink(file), '/'):
          if options.showOk:
            echoStatus("OK/", file)
        else:
          if options.showOk:
            echoStatus("OK", file)
      else:
        echoStatus("ERR_SYM_FILE", file)
        # possibly fixable, see reasoning in runSymlinkDir()
        echoPoint("(possibly fixable)")
    # symlink created by user
    else:
      if options.showOk:
        echoStatus("OK_USYM_FILE", file)

  proc runSymlinkDir(file: string, real: string, options: Options): void =
    if symlinkCreatedByDotty(dotDir, homeDir, real):
      if symlinkResolvedProperly(dotDir, homeDir, file):
        if endsWith(expandSymlink(file), '/'):
          if options.showOk:
            echoStatus("OK/", file)
        else:
          if options.showOk:
            echoStatus("OK", file)
      else:
        echoStatus("ERR_SYM_DIR", file)

        # possibly fixable because when we have this:
        # ~/.profile (file) -> ~/.dots/.config/profile/.profile (real)
        # it becomes this:
        # ~/.profile (file) -> ~/.dots/.profile (real)
        # even though, it should be
        # ~/.profile (file) -> ~/.config/profile/.profile (real)
        # this a user error symlinking inside dotDir, but nevertheless,
        # still not necessarily fixable. we can't fix this
        echoPoint("(possibly fixable)")
    # symlink created by user
    else:
      if options.showOk:
        echoStatus("OK_USYM_DIR", file)

  proc runSymlinkNull(file: string, real: string): void =
    if symlinkCreatedByDotty(dotDir, homeDir, real):
      echoStatus("ERR_SYM_NULL", file)
      echoPoint(fmt"{file} (symlink)")
      echoPoint(fmt"{real} (nothing here)")
      echoPoint(fmt"Did you forget to create your actual dotfile at '{real}'?")
      echoPoint("(not fixable)")
    # symlink created by user
    else:
      echoStatus("ERR_USER_SYM", file)
      echoPoint(fmt"{file} (symlink)")
      echoPoint(fmt"{real} (nothing here)")
      echoPoint(fmt"Did you forget to create your actual dotfile at '{real}'?")
      echoPoint("(not fixable)")

  proc runFileFile(file: string, real: string): void =
    echoStatus("ERR_FILE_FILE", file)
    echoPoint(fmt"{file} (file)")
    echoPoint(fmt"{real} (file)")
    echoPoint("(possibly fixable)")

  proc runFileDir(file: string, real: string): void =
    echoStatus("ERR_FILE_DIR", file)
    echoPoint(fmt"{file} (file)")
    echoPoint(fmt"{real} (directory)")
    echoPoint("(not fixable)")

  proc runFileNull(file: string, real: string): void =
    echoStatus("ERR_FILE_NULL", file)
    echoPoint("(fixable)")

  proc runDirFile(file: string, real: string): void =
    echoStatus("ERR_DIR_FILE", file)
    echoPoint(fmt"{file} (directory)")
    echoPoint(fmt"{real} (file)")
    echoPoint("(not fixable)")

  proc runDirDir(file: string, real: string): void =
    echoStatus("ERR_DIR_DIR", file)
    echoPoint(fmt"{file} (directory)")
    echoPoint(fmt"{real} (directory)")
    echoPoint("Remove the directory that has the older contents")
    echoPoint("(possibly fixable)")

  proc runDirNull(file: string, real: string): void =
    echoStatus("ERR_DIR_NULL", file)
    echoPoint("(fixable)")

  proc runNullFile(file: string, real: string): void =
    echoStatus("ERR_NULL_FILE", file)
    echoPoint("(fixable)")

  proc runNullDir(file: string, real: string): void =
    echoStatus("ERR_NULL_DIR", file)
    echoPoint("(fixable)")

  proc runNullNull(file: string, real: string): void =
    echoStatus("ERR_NULL_NULL", file)
    echoPoint(fmt"Did you forget to create your actual dotfile at '{real}'?")
    echoPoint("(not fixable)")

  doAbstract(
    dotDir,
    homeDir,
    options,
    dotFiles,
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
    dotFiles: seq[string]) =
  proc runSymlinkSymlink(file: string, real: string, options: Options): void =
    if symlinkCreatedByDotty(dotDir, homeDir, real):
      echoStatus("ERR_SYM_SYM", file)
      echoPoint("(not fixable)")

  proc runSymlinkFile(file: string, real: string, options: Options) =
    if symlinkCreatedByDotty(dotDir, homeDir, real):
      if symlinkResolvedProperly(dotDir, homeDir, file):
        # if destination has an extraneous forward slash,
        # automatically remove it
        if endsWith(expandSymlink(file), '/'):
          let temp = expandSymlink(file)
          removeFile(file)
          createSymlink(rts(temp), file)
      else:
        echoStatus("ERR_SYM_FILE", file)
        echoPoint("(attempted fix)")

        # won't always work properly see reasoning in status runSymlinkDir()
        removeFile(file)
        createSymlink(getRealDot(dotDir, homeDir, file), file)

  proc runSymlinkDir(file: string, real: string, options: Options) =
    if symlinkCreatedByDotty(dotDir, homeDir, real):
      if symlinkResolvedProperly(dotDir, homeDir, file):
        # if destination has a spurious slash, automatically
        # remove it
        if endsWith(expandSymlink(file), '/'):
          let temp = expandSymlink(file)
          removeFile(file)
          createSymlink(rts(temp), file)
      else:
        # won't always work properly see reasoning in status runSymlinkDir()
        removeFile(file)
        createSymlink(getRealDot(dotDir, homeDir, file), file)

  proc runSymlinkNull(file: string, real: string) =
    if symlinkCreatedByDotty(dotDir, homeDir, real):
      echoStatus("ERR_SYM_NULL", file)
      echoPoint(fmt"{file} (symlink)")
      echoPoint(fmt"{real} (nothing here)")
      echoPoint("(not fixable)")

  proc runFileFile(file: string, real: string) =
    let fileContents = readFile(file)
    let realContents = readFile(real)

    if fileContents == realContents:
      removeFile(file)
      createSymlink(real, file)
    else:
      echoStatus("ERR_FILE_FILE", file)
      echoPoint(fmt"{file} (file)")
      echoPoint(fmt"{real} (file)")
      echoPoint("(not fixable)")

  proc runFileDir(file: string, real: string) =
    echoStatus("ERR_FILE_DIR", file)
    echoPoint(fmt"{file} (file)")
    echoPoint(fmt"{real} (directory)")
    echoPoint("(not fixable)")

  proc runFileNull (file: string, real: string) =
    echoStatus("ERR_FILE_NULL", file)
    echoPoint("Automatically fixed")

    createDir(parentDir(real))

    # file doesn't exist on other side. move it
    moveFile(file, real)
    createSymlink(real, file)

  proc runDirFile (file: string, real: string) =
    echoStatus("ERR_DIR_FILE", file)
    echoPoint(fmt"{file} (directory)")
    echoPoint(fmt"{real} (file)")

  # swapped
  proc runDirNull (file: string, real: string) =
    # ensure directory
    createDir(parentDir(real))

    # file doesn't exist on other side. move it
    try:
      echoStatus("ERR_DIR_NULL", file)
      echoPoint("Automatically fixed")

      copyDirWithPermissions(file, real)
      removeDir(file)
      createSymlink(real, file)
    except Exception:
      logError getCurrentExceptionMsg()
      echoStatus("ERR_DIR_NULL", file)
      echoPoint("Error: Could not copy folder")

  # swapped
  proc runDirDir (file: string, real: string) =
    if dirLength(file) == 0:
      echoStatus("ERR_DIR_DIR", file)
      echoPoint("Automatically fixed")

      removeDir(file)
      createSymlink(getRealDot(dotDir, homeDir, file), file)
    elif dirLength(real) == 0:
      echoStatus("ERR_DIR_DIR", file)
      echoPoint("Automatically fixed")

      removeDir(real)
      runDirNull(file, real)
    else:
      echoStatus("ERR_DIR_DIR", file)
      echoPoint(fmt"{file} (directory)")
      echoPoint(fmt"{file} (directory)")
      echoPoint("(not fixable)")

  proc runNullAny(file: string, real: string) =
    createSymlink(getRealDot(dotDir, homeDir, file), file)

  doAbstract(
    dotDir,
    homeDir,
    options,
    dotFiles,
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
    runNullAny,
    runNullAny,
    runNullAny
  )
