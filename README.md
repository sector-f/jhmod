# jh_extract

A tool to work with Jupiter Hell game files.

## Install

### Manually

Refer to https://go.dev/doc/tutorial/compile-install for more information on
how to install golang packages manually.

```bash
# For those who don't have it already defined.
export GOPATH="$HOME/go"
export PATH="$PATH:$GOPATH/bin"

go install github.com/sector-f/jh_extract
```

### Nix

TODO

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

#### 2) build the pathlist using `jh_extract pathlist scan`

```bash
$ jh_extract pathlist scan core.324962 | tee pathlist.txt
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
jh_extract knows about all of the paths.  Using an old pathlist.txt may result
in many unknown paths that need to be extracted with `-u`.**

```bash
$ jh_extract extract -f core.nvc -p samples/pathlist.txt
Extracted 266 files
$ ls data/
font  lang  lua
```

#### Q: Where are my assets?

You may need a more complete pathlist.txt.  The one provided is reproducible by
using `jh_extract pathlist scan`, however, it doesn't find graphics data
(TODO).

## Jupiter Hell Modding cheatsheet

1. Find your game directory that contains `.nvc` files.  There should be
   `assets.nvc` and `core.nvc`.  Note down these paths.
2. Use `jh_extract extract -f core.nvc -p samples/pathlist.txt` to extract the
   files into the current directory.  Most (all?) files should be located in a
   `data` subdirectory.
3. Build an "overlay" directory that mirrors the directory structure of the
   extracted paths.  It only has to have the directories needed to place your
   modified game files.
4. Put your mod files in this overlay directory.
5. Copy the overlay directory tree into your JH gamedir such that your overlay
   directory has its `data/` directory adjacent to the `core.nvc` file.
6. Viola you can test your mod.

## License

BSD 3-Clause ([BSD-3-Clause](https://spdx.org/licenses/BSD-3-Clause.html)).
See [LICENSE](./LICENSE).
