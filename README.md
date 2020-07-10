# Globe

A flexible and language-agnostic configuration system for managing multiple independent projects. It automatically manages the modification or bootstrapping of files that are nearly always used in development workflows or project repositories. This is most useful when develvoping with many tens or hundreds of projects, especially if they're hosted on GitHub

## Table of Contents

## Package Contents

### `tools/bootstrap`

Utilities for bootstrapping new projects

### `tools/fix`

Utilities for autofixing random files

## Configuration that's Automatically Managed

-  GitHub
-  EditorConfig
-  Git
-  License
-  common scripts

'common scripts' are scripts I use across my projects. they are local to projects rather than my system to increase code modularity and decrease errors that may spawn from code collaboration and execution across independent or inconsistent environments

### Limitations

Some files are long and aren't templated to only have the necessary lines (see issue #2 for details)

## Usage

## Installation

```sh
git clone https://github.com/eankeen/globe
go install
```

## License
