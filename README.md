# Dotty

🌎 Ensure common configuration across projects

## Description

A CM (Configuration Management) utility for dotfiles. It's used for managing local, user, or system-wide dotfiles.

For example, you can manage your "local" `.editorconfig`'s, `.eslintrc.js`'s, `.clang-format`'s, your "user" `~/.bashrc`'s, `~/.inpurc`'s, or your "global" `/boot/efi/EFI/refind/refind.conf`, `/root/.nanorc` files.

## Usage

```txt
$ dotty --help
A CM (Configuration Management) utility for dotfiles. Used for managing local, user, or system-wide dotfiles

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
  -h, --help             help for globe

Use "globe [command] --help" for more information about a command.
```

## Installation

```sh
git clone https://github.com/eankeen/dotty
go install
```
