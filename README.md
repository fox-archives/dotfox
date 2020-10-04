# Dotty

🌎 Ensure common configuration across projects

## Description

A CM (Configuration Management) utility for dotfiles. It's used for managing local, user, or system-wide dotfiles. For example, you can manage your "local" `.editorconfig`'s, `.eslintrc.js`'s, `.clang-format`'s, your "user" `~/.bashrc`'s, `.inpurc`'s, or your "global" `/boot/efi/EFI/refind/refind.conf`, `/root/.nanorc` files.

## Usage

```txt
$ globe --help
Command:
  globe

Description:
  An easy to use language-agnostic configuration management tool

Commands:
  init    Initiate Globe configuration
  sync    Update configuration and files

Options:
  --help Display help menu
```

## Installation

```sh
git clone https://github.com/eankeen/globe
go install
```
