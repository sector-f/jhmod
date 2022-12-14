package gui

import (
	"runtime"

	"github.com/adrg/xdg"
)

func guessGameDir() string {
	os := runtime.GOOS
	switch os {
	case "windows":
		return "C:/Program Files (x86)/Steam/steamapps/common/Jupiter Hell"
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
