# Usage

Before I move the docs to a cooler place, here are the features currently available.

You can manage three types of dotfiles. _system_, _user_, and _local_. When setting up your repository, make sure you have a `config`, `local`, `system`, and `local` folders [like mine for example](https://github.com/eankeen/dots). I have my repo cloned to `~/.dots` and have this alias setup

```sh
alias dotty='dotty --dotfiles-dir=$HOME/.dots'
```

## Types of Dotfiles

### System

ALPHA

System isn't quite done yet - the code doesn't work. Don't worry about accidentally using it because you need to have elevated privileges to use it

### User

The config format looks like this:

```toml
#
# ignores
#

[[ignores]]
file = "/bash-it/"

[[ignores]]
file = "/oh-my-zsh/"

#
# ~
#

[[files]]
file = "/.bash_logout"
tags = [
	"bash"
]

[[files]]
file = "/.bash_profile"
tags = [
	"bash"
]
```

#### Ignore Entry

- specify a file (required)
- tags don't do anything yet
- for each ignore entry, the trailing `/` is removed
  - this _might_ change in the future to do more intelligent matching, not sure
- in the end, a simple `strings.Contains()` (substring match) is performed to make sure it matches
  - for example, if you have an entry called `file = "oh-my-zsh"`, and you have a folder at `/home/edwin/.config/zsh/oh-my-zsh`, naturally everything in the folder (including the folder itself) will be excluded
- this is useful because `oh-my-zsh` has a child folder called `scripts`, but i don't want it to interfere with _my_ `scripts` folder

#### File Entry

- specify a file (required)
- tags don't do anything yet
- **important feature**
  - if you have a trailing slash, it will only select a matching folder
  - if you don't have a trailing slash, it will only select a matching file
- in the end, a `strings.HasSuffix()` check is performed to make sure it matches
  - fir example, if you have a file at `/home/edwin/.config/vim/vimrc`, and you have a `file = "vimrc"`, then it will check to see if `.config/vim/vimrc` has a suffix of `vimrc`) (yes, the `/home/edwin` is sliced out)

### Local

The previous implementation was taken from another project of mine - all of the code was ripped out, all it does is write content to a specific file with some hardcoded values. It's not ready to be used.

## Useful Commands

## apply

You can sync the dotfiles with:

```sh
dotty <system|user|local> apply
```

## edit

You can edit your dotfiles with:

```sh
dotty <system|user|local> edit
```

It first looks for your EDITOR environment variable and fallbacks to `vim` as the editor

## check

ALPHA

You can check to see if your config contains ambiguities or if some values aren't used. (Useful because you only specify the minimum required path in the config, so you may not actually be matching a particular file)

```sh
dotty user check
```

NOTE:

there are currently three heuristics for determining if a dotfile entry (ex. `file = "/.vimrc"`) is valid. only the first two work as intended, but take it with a grain of salt because it's not finished
