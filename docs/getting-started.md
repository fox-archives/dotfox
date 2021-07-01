# Getting Started

When setting up Dotty, you need two things: a dot directory (`dotDir`) and a destination directory (`destDir`). `dotDir` is a directory that your dotfiles tracked with your VCS of choice, such as `~/.dotfiles`. `destDir` is the location to deploy your dotfiles to, such as `~/`.

Specify these two things in `~/.config/dotty/config.toml` (or with using the config directory of your choice with `XDG_CONFIG_HOME`)

```toml
[config]
dotDir = "~/.dotfiles"
destDir = "~"
```

Now, specify the dotfiles you wish to automatically deploy. These dotfiles will have symlinks created in the `destDir`, pointing to their respective file or directory in `dotDir`. Specify the dotfiles using a shell script:

```bash
#!/usr/bin/env bash

declare -r home="$HOME"
declare -r cfg="${XDG_CONFIG_HOME:-$HOME/.config}"
declare -r data="${XDG_DATA_HOME:-$HOME/.local/share}"

declare -ra dotfiles=(
	"$home/.bashrc"
	"$home/.bash_logout"
)

for dotfile in "${dotfiles[@]}"; do
	printf "%s\n" "$dotfile"
done
```

Dotty will execute this script, and use every line of standad output as a separate dotfile to track. Standard output would look like the following in this case

```txt
/home/edwin/.bashrc
/home/edwin/.bash_logout
```

Note that with every line, there is always a prefix of the `destDir` (`/home/edwin` i.e. `~/`)

Now, let's try running Dotty

```sh
$ dotty status
[ERR_NULL_NULL] /home/edwin/.bashrc
                -> (not fixable)
                -> Is there a file or directory at /home/edwin/.dotfiles/.bashrc?
[ERR_NULL_NULL] /home/edwin/.bash_logout
                -> (not fixable)
                -> Is there a file or directory at /home/edwin/.dotfiles/.bash_logout?
Done.
```

As you can see, there is a hint that tells us that we forgot to move our dotfiles to the `dotDir`. How useful! Let's do that...

```sh
echo 'export EDITOR="vim"' > ~/.dotfiles/.bashrc
echo 'clear' > ~/.dotfiles/.bash_logout
```

Now, let's try running again

```sh
$ dotty status
[ERR_NULL_FILE] /home/edwin/.bashrc
                -> (fixable)
[ERR_NULL_FILE] /home/edwin/.bash_logout
                -> (fixable)
Done.
```

As you can see, the `ERR_NULL_NULL` status codes changed to `ERR_NULL_FILE`. The last `NULL` changed to a file because we placed the correct files in the `dotDir` (destination directory)

Now that the dotfiles are automatically fixable, let's run `dotty reconcile`

```sh
$ dotty reconcile
Done.
$ dotty status
[OK]            /home/edwin/.bashrc
[OK]            /home/edwin/.bash_logout
Done.
```

Cool, now the dotfiles have been resolved correctly. If you really want, you can check to see if the symlinks point to the correct location

```sh
$ ls -al ~/.bash*
lrwxrwxrwx 1 edwin edwin   34 Jun 29 21:36 /home/edwin/.bash_logout -> /home/edwin/.dotfiles/.bash_logout
lrwxrwxrwx 1 edwin edwin   29 Jun 29 21:36 /home/edwin/.bashrc -> /home/edwin/.dotfiles/.bashrc
```

And they do have their correct target files! 