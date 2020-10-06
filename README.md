# Dotty

🌎 System, user, and local specific dotfile manager

## Description

A CM (Configuration Management) utility for dotfiles. It's used for managing local, user, or system-wide dotfiles.

For example, you can manage your "local" `.editorconfig`'s, `.eslintrc.js`'s, `.clang-format`'s, your "user" `~/.bashrc`'s, `~/.inpurc`'s, or your "system" `/boot/efi/EFI/refind/refind.conf`, `/root/.nanorc` files.

## Features

- Human readable config format
- Prompts user on conflicts
- Uses symlinks
- Use the same utility to manage three different types of dotfiles

## Usage

```txt
$ dotty --help
A CM (Configuration Management) utility for dotfiles. Used for managing
local,user, or system-wide dotfiles

Usage:
   dotty [command]

Available Commands:
  help        Help about any command
  init        Init Dotty's config files
  local       Local (.) (per-project) config management
  system      Systemwide (/) config management
  user        Userwide (~) config management

Flags:
      --dot-dir string   The location of your dotfiles
  -h, --help             help for dotty

Use "dotty [command] --help" for more information about a command.
```

## Installation

```sh
git clone https://github.com/eankeen/dotty
go install
```
