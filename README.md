# jhmod

[![Build](https://github.com/sector-f/jhmod/actions/workflows/build.yml/badge.svg)](https://github.com/sector-f/jhmod/actions/workflows/build.yml)


A tool to work with Jupiter Hell game files.

## Features

- Create and extract `.nvc` archives
- Scan for interesting `.nvc` archive paths referenced in the JH program
- Get information from save files


## Install

### Manually

Refer to https://go.dev/doc/tutorial/compile-install for more information on
how to install golang packages manually.

```bash
# For those who don't have it already defined.
export GOPATH="$HOME/go"
export PATH="$PATH:$GOPATH/bin"

go install github.com/sector-f/jhmod
```

### Nix

A flake is offered in this repository.  You can add it to your own flake as a
input, then add it to your `environment.systemPackages`.  See example usage
[here](https://gitlab.com/search?project_id=36950231&search=jhmod).

You can develop using the flake via `nix develop` or setting up direnv so it
auto-loads the flake.  Add `use flake` to your `.envrc` then run `direnv
allow`.


## Usage

### Conventions

- Lines beginning with `$` are commends you type.
- Lines beginining with `... SNIP ...` indicate lines that were removed for
  brevity.
- Other lines are program output.

### Build a list of interesting paths

#### 1) obtain a core file using `./scripts/get-core`.

```bash
$ ./scripts/get-core
core.324962
```

#### 2) build the pathlist using `jhmod pathlist scan`

```bash
$ jhmod pathlist scan core.324962 | tee pathlist.txt
data/font/terminal8x9_c64i.png
data/lang/de.csv
data/lang/pl.csv
... SNIP ...
data/lua/jh/gfx/tilesets/ts01/ts01.lua
data/lua/jh/gfx/tilesets/ts01/ts01_A.lua
data/lua/jh/gfx/tilesets/ts01/ts01_B.lua
```

### Extract `.nvc` files

**Note: You should consider generating your own pathlist.txt to ensure
jhmod knows about all of the paths.  Using an old pathlist.txt may result
in many unknown paths that need to be extracted with `-u`.**

```bash
$ jhmod extract -f core.nvc -p samples/pathlist.txt
Extracted 266 files
$ ls data/
font  lang  lua
```

#### Q: Where are my assets?

You may need a more complete pathlist.txt.  The one provided is reproducible by
using `jhmod pathlist scan`, however, it doesn't find graphics data
(TODO).

## Jupiter Hell Modding cheatsheet

1. Find your game directory that contains `.nvc` files.  There should be
   `assets.nvc` and `core.nvc`.  Note down these paths.
2. Use `jhmod extract -f core.nvc -p samples/pathlist.txt` to extract the
   files into the current directory.  Most (all?) files should be located in a
   `data` subdirectory.
3. Create a folder in your JH game directory called `mods`
4. Create a subfolder in your JH game directory that will contain your mod, say `coolestmod`
5. Put some mod lua code into `mods/coolestmod/main.lua`.  This will run when the game loads.
6. Use the files created via `jhmod extract ...` to determine what lua tables
   to manipulate to achieve your desired effect.
7. You can now test your mod.

See [Modding](https://jupiterhell.fandom.com/wiki/Modding) on the Jupiter Hell
Wiki for more information.

## License

BSD 3-Clause ([BSD-3-Clause](https://spdx.org/licenses/BSD-3-Clause.html)).
See [LICENSE](./LICENSE).
