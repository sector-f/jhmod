package gui

import (
	"os"
	"runtime"

	path "path/filepath"

	"github.com/adrg/xdg"
)

const (
	GAMEDIR_BASE = "Jupiter Hell"
	COMMONDIR    = "Steam/steamapps/common"
)

func guessGameDir() string {
	currentOS := runtime.GOOS
	switch currentOS {
	case "windows":
		return path.Join("C:/Program Files (x86)", COMMONDIR, GAMEDIR_BASE)
	case "linux":
		dirname, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		return path.Join(dirname, ".local/share", COMMONDIR, GAMEDIR_BASE)
	default:
		panic("Don't know about OS default game dir")
	}
}

func getSaveScumDir() string {
	dir, err := xdg.DataFile("jhmod/savescum")
	if err != nil {
		panic(err)
	}
	return dir
}
