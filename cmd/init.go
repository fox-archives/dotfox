package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/eankeen/dotty/internal/util"
	"github.com/eankeen/go-logger"
	"github.com/spf13/cobra"
)

func dir(path string) {
	wd, err := os.Getwd()
	util.HandleFsError(err)
	err = os.MkdirAll(filepath.Join(wd, path), 0755)
	util.HandleFsError(err)
}

func file(path string, content string) {
	wd, err := os.Getwd()
	util.HandleFsError(err)
	err = ioutil.WriteFile(filepath.Join(wd, path), []byte(content), 0644)
	util.HandleFsError(err)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Dotfiles repository to be managed with Dotty",
	Run: func(cmd *cobra.Command, args []string) {
		// root
		file("dotty.toml", `configDir = "cfg"

# relative to root of project (containing dotty.toml)
systemDirSrc = "system"
systemDirDest = "/"

userDirSrc = "user"
userDirDest = "~"

localDirSrc = "local"
`)

		// config
		dir("cfg")
		file("cfg/system.dots.toml", `[[files]]
file = "/etc/systemd/journald.conf
tags = [ "systemd" ]

[[files]]
file = "/etc/binfmt.d/10-go.conf
`)
		file("cfg/user.dots.toml", `[[files]]
file = "/.profile"
tags = [ "shell" ]

[[files]]
file = "/.bashrc"
tags = [ "bash" ]

[[files]]
file = "/.bash_logout"
tags = [ "bash" ]
`)
		file("cfg/local.dots/toml", `[[files]]
		file = "/.editorconfig"

		[[files]]
		file = "/.prettierrc.json"
`)

		// system
		dir("system")
		dir("system/etc/systemd")
		file("system/etc/systemd/journald.conf", `#  This file is part of systemd.
#
# See journald.conf(5) for details.

[Journal]
Storage=auto
Compress=yes
`)

		dir("system/etc/binfmt.d")
		file("system/etc/binfmt.d/10-go.conf", `:golang:E::go::/usr/bin/gorun:OC
`)

		// user
		dir("user")
		file("user/.profile", `#
# ~/.profile
#

export LANG="en_US.UTF-8"
export VISUAL="vim"
export PAGER="less"
export CMD_ENV="linux"
`)
		file("user/.bashrc", `#
# ~/.bashrc
#

[ -r ~/.profile ] && source ~/.profile

[[ $- != *i* ]] && return

HISTCONTROL="ignorespace:ignoredups"
HISTIGNORE="?:ls:[bf]g:pwd:clear*:exit*"
HISTTIMEFORMAT="%B %m %Y %T | "

shopt -s autocd
shopt -s cdable_vars
shopt -s checkjobs
`)
		file("user/.bash_logout", `#
# ~/.bash_logout
#

clear; printf '\033[3J'
`)

		// local
		dir("local")
		file("local/.editorconfig", `root = true

[*]
indent_style = tab
indent_size = unset
tab_width = 4
end_of_line = lf
charset = utf-8
trim_trailing_whitespace = true
insert_final_newline = true

[{GNUmakefile,Makefile,makefile,*.mk}]
indent_style = tab

[*.{yaml,yml}]
indent_style = space

[*.md]
indent_style = space
`)
		file("local/.prettierrc.json", `{
	"semi": false,
	"singleQuote": true
}`)

		logger.Informational("Successfully initiated dotfile repository\n")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
