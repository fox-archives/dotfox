# Globe

Language-agnostic configuration and program management, and build utility

## Description

A language-agnostic generalized programming utility for managing multiple independent projects. It automatically manages the bootstraping and synchronizing of files in addition to the bootstrapping, synchronizing, and installation of programs that are repetitevly used in development workflows or project repositories. This is intended to significantly decrease miscellaneous (yet somewhat common) development tasks when working across many tens or hundreds of projects

![xkcd](https://imgs.xkcd.com/comics/is_it_worth_the_time.png)

Right now, it's not extendable and specific to my situation. I'm working on generalizing this utility so it can be glued together with other independent (but related) tools

## Table of Contents

## Configuration that's Automatically Managed

-  GitHub
-  EditorConfig
-  Git
-  License
-  common scripts (`tools/`)

'common scripts' are scripts I use across my projects. they are local to projects rather than my system to increase code modularity and decrease errors that may spawn from code collaboration and execution across independent or inconsistent environments. Right now I'm trying to find the right balance of these little scripts being located in `tools` or as builtin to the main Go package. Maybe they shouldn't exist here at all and tools should live in separate locations and their execution should be delegated from a task runner

### Limitations

Some files are long and aren't templated to only have the necessary lines (see issue #2 for details)

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
