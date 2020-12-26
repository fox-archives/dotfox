import os
import system
import strutils
import logging

import "./util"

let consoleLog = newConsoleLogger()
addHandler(consoleLog)

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

proc doReconcile(dotDir: string, homeDir: string, dotConfigFiles: seq[string]) =
  for i, file in dotConfigFiles:
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

let home = getHomeDir() & "/"
let cfg = xdgCfg() & "/"
let data = xdgData() & "/"
let dotConfigFiles = [
  home & ".alsoftrc",
  home & ".bash_logout",
  home & ".bash_profile",
  home & ".bashrc",
  home & ".cliflix.json",
  home & ".profile",
  home & ".yarnrc",
  cfg & "alacritty",
  cfg & "awesome",
  cfg & "bash",
  cfg & "bat",
  cfg & "broot",
  cfg & "calcuse",
  cfg & "cava",
  cfg & "chezmoi",
  cfg & "cmus/rc",
  cfg & "Code/User/keybindings.json",
  cfg & "Code/User/settings.json",
  cfg & "curl",
  cfg & "dircolors",
  cfg & "dunst",
  cfg & "environment.d",
  cfg & "fish",
  cfg & "fontconfig",
  cfg & "gh",
  cfg & "git",
  cfg & "hg",
  cfg & "htop",
  cfg & "i3",
  cfg & "i3status",
  cfg & "info",
  cfg & "ion",
  cfg & "kitty",
  cfg & "lazydocker",
  cfg & "maven",
  cfg & "micro",
  cfg & "mnemosyne/config.py",
  cfg & "nano",
  cfg & "neofetch",
  cfg & "nimble",
  cfg & "nitrogen",
  cfg & "npm",
  cfg & "nu",
  cfg & "nvim",
  cfg & "pamix.conf",
  cfg & "picom",
  cfg & "polybar",
  cfg & "profile",
  cfg & "pulse/client.conf",
  cfg & "pulsemixer.cfg",
  cfg & "python",
  cfg & "ranger",
  cfg & "readline",
  cfg & "redshift",
  cfg & "ripgrep",
  cfg & "rofi",
  cfg & "rtorrent",
  cfg & "salamis",
  cfg & "starship",
  cfg & "sxhkd",
  cfg & "systemd",
  cfg & "terminator",
  cfg & "termite",
  cfg & "tilda",
  cfg & "tmux",
  cfg & "todotxt",
  cfg & "urxvt",
  cfg & "user-dirs.dirs",
  cfg & "vim",
  cfg & "wget",
  cfg & "X11",
  cfg & "xob",
  cfg & "yay",
  cfg & "zsh",
  data & "gnupg/gpg.conf",
  data & "gnupg/dirmngr.conf",
  data & "applications/Calcurse.desktop",
  data & "applications/Nemo.desktop",
  data & "applications/Obsidian.desktop",
  data & "applications/Zettlr.desktop"
]

let homeDir = getHomeDir()
let dotDir = joinPath(getHomeDir(), ".dots/user")

if paramCount() < 1:
  echo "Error: Expected subcommand. Exiting"
  quit 1

case paramStr(1):
  of "status":
    doStatus(dotDir, homeDir, @dotConfigFiles)
  of "reconcile":
    doReconcile(dotDir, homeDir, @dotConfigFiles)
  else:
    echo "Error: Subcommand not found. Exiting"
    quit 1
