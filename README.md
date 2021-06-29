# Dotty

🌎 Simple dotfile manager

## Description

Simple, clean tool to automatically track your dotfiles and ensure they are all resolved properly. With the files and directories you specify, it automatically creates symlinks to them in the destination (home) directory

## Setup and Configuration

Dotty will read the following base config folder

```toml
# $XDG_CONFIG_HOME/dotty/config.toml
[config]
dotDir = "~/.dots/user"
destDir = "~"
```

In this case, I have dotfiles like `.bashrc`, `.bash_logout` located in `~/.dots/user`

Specify these files in `$XDG_CONFIG_HOME/dotty/dotty.sh`

```bash
#!/usr/bin/env bash

home="$HOME"
cfg="${XDG_CONFIG_HOME:-$HOME/.config}"
data="${XDG_DATA_HOME:-$HOME/.local/share}"

declare -ra dotFiles=(
	"$home/.bashrc"
	"$home/.bash_logout"
)

for dotFile in "${dotFiles[@]}"; do
	printf "%s\n" "$dotFile"
done
```

Dotty reads each dotfile separated by newline, so a simpler example could be:

```sh
#!/usr/bin/env sh
echo "/home/$(whoami)/.bashrc"
echo "/home/$(whoami)/.bash_logout
```

Note that each of these entries have a prefix that is the same of `destDir` in the `config.toml` file

Next, run

```sh
dotty status
```

Read the [Status Codes](##Status Codes) section below to better interpret the output

Run `dotty reconcile` to automatically setup your dotfiles. In doing this, it _always_ creates a symlink in the `destDir` for each file you specify in `dotty.sh` to its respective location in `dotDir` (ex. `~/.dots`). If symlinks don't already exists, it copies your dots to `dotDir` and puts symlinks in their place; automatically resolving as many dotfiles as possible

## Status Codes

```txt
FORMAT
  [generalStatus]-[homeDirStatus]-[dotDirStatus]

generalStatus:
  OK:
    All symlinks and files exists properly

  OK/:
    All symlinks and files exist properly. And, the symlink destination
    has an extraneous succeeding slash ("even when the destination is a
    directory this is not needed"). So on reconcile this will normalize
    to be suffix-slash-less.

  ERR:
    A conflict exists. You yourself may need to delete one of the files if dotty
    is not able to autofix. When autofixing, the tool will ensure the
    non-deleted file is at ~/.dots/.profile1, and a symlink to it exists
    in ~/.profile1


homeDirStatus:
  SYM
    Symlink exists

  FILE
    File exists

  DIR
    Directory exists

  NULL
    Nothing exists

  USYM
    Symlink exists, but it's not managed by 'dotty'. Occurs when
    a traced symlink's destination is not within dotDir
    (ex. ~/.profile -> ~/.config/profile/profile.sh). All lines with USYM
    are OK because they are symlinks this tool does not manage

dotDirStatus:
  SYM
    Symlink exists

  FILE
    File exists

  DIR
    Directory exists

  NULL
    Nothing exists

EXAMPLES
  ERR_NULL_FILE
  OK_USYM_DIR

```
