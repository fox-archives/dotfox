# Reference

## Status Codes

The status code is in the format of

```txt
<generalCode>-<destDirCode>-<dotDirCode>
```

### `generalCode`

The general status code shows the general validity of the symlink (whether it points to the correct location)

#### `OK`

The symlink correctly points to the target file or directory

#### `OK/`

The symlink correctly points to the target file or directory, with the caveat that the target name has an extraneous suffix slash. On `dotfox deploy`, this will automatically be normalized

#### `ERR`

A conflict exists. DotFox will attempt to automatically fix the problem. However, if dotfox is unable to, you will need to fix it yourself

### `destDirCode`

The `destDirCode` represents status code of a particular dotfile with respect to its destination directory. The destination directory is the directory in which your dotfiles are deployed to and the directory the symlinks are created. Usually, this is located at `~/` For a healthy deployment, you would want these to have a value of `SYM`

#### `SYM`

The particular dotfile has a corresponding symlink in the destination directory

#### `FILE`

The particular dotfile has a corresponding file in the destination directory. DotFox will try to automatically replace this with a symlink pointing to said target dotfile on deploy

#### `DIR`

The particular dotfile has a corresponding directory in the destination directory. DotFox will try to automatically replace this with a symlink pointing to said target dotfile on deploy

#### `NULL`

The particular dotfile has no corresponding file or directory in the destination directory. DotFox will try to automatically replace this with a symlink pointing to said target dotfile on deploy

### `dotDirCode`

The `dotDirCode` represents status code of a particular dotfile with respect to its dotfile directory. The destination directory is the directory in which your dotfiles are held in version control. Most people choose to place this directory at `~/.dotfiles` or `~/.dots`. For a healthy deployment, you would want these to have a value of either `FILE`, `DIR`, or `SYM`

### `SYM`

The particular dotfile has a corresponding symlink in the dot directory. This is possible, for example, if you wish to symlink a file like `~/.bashrc` to `~/.config/bash/bashrc.sh`

### `FILE`

The particular dotfile has a corresponding file in the dot directory

### `DIR`

The particular dotfile has a corresponding directory in the dot directory

### `NULL`

The particular dotfile has a no corresponding file or directory in the dot directory.
