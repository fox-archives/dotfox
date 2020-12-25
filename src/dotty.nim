#!/usr/bin/env nimcr
import os
import system
import strutils
import logging

proc xdgCfg: string =
   if getEnv("XDG_CONFIG_HOME") != "":
      return getEnv("XDG_CONFIG_HOME")
   else:
      return joinPath(getEnv("HOME"), ".config")

proc xdgData: string =
   if getEnv("XDG_DATA_HOME") != "":
      return getEnv("XDG_DATA_HOME")
   else:
      return joinPath(getEnv("HOME"), ".local/share")

let home = getEnv("HOME") & "/"
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

# remove trailing slash
proc rts(str: string): string =
   if endsWith(str, '/'):
      return str[0 .. ^2]
   return str



proc symlinkResolvedProperly(src: string, dest: string, symlink: string): bool =
   let rel = symlink[len(dest) .. ^1]

   if rts(expandSymlink(symlink)) == joinPath(src, rel):
      return true
   else:
      return false

# checks to see if src and dest are both a folder or file. if not, aborts
proc isBothFileOrFolder(src: string, dest: string, rel: string): bool =
   let pathOne = joinPath(src, rel)
   let pathTwo = joinPath(dest, rel)

   if dirExists(pathOne) and dirExists(pathTwo):
      return true
   elif fileExists(pathOne) and fileExists(pathTwo):
      return true
   else:
      return false


proc status(src: string, dest: string) =
   for i, file in dotConfigFiles:
      if symlinkExists(file):
         if symlinkResolvedProperly(src, dest, file):
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
         let rel = file[len(dest) .. ^1]
         if fileExists(joinPath(src, rel)):
            echo "[TOLINK_F] " & file
         elif dirExists(joinPath(src, rel)):
            echo "[TOLINK_D] " & file
         else:
            echo "[MISSING]  " & file

proc reconcile(src: string, dest: string) =
   for i, file in dotConfigFiles:
      if symlinkExists(file):
         if symlinkResolvedProperly(src, dest, file):
            # everything OK except the trailing slash; fix it
            if endsWith(expandSymlink(file), '/'):
               let temp = expandSymlink(file)
               removeFile(file)
               createSymlink(rts(temp), file)
         else:
            let temp = expandSymlink(file)
            removeFile(file)
            createSymlink(rts(temp), file)
      else:
         let rel = file[len(dest) .. ^1]
         createSymlink(joinPath(src, rel), file)

let dest = getEnv("HOME")
let src = joinPath(getEnv("HOME"), ".dots/user")

status(src, dest)
# reconcile(srf, dest)
