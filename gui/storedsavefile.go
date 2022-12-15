package gui

import (
	"os"
	path "path/filepath"

	"github.com/u-root/u-root/pkg/cp"
	"gorm.io/gorm"
)

type StoredSaveFile struct {
	gorm.Model
	Id            int `gorm:"primary_key"`
	OriginalBase  string
	StoredRelPath string `gorm:"unique"`
	Sha256Hex     string `gorm:"unique"`

	// Copied verbatim from savefile.go
	// Player name
	PlayerName string
	// Game mode.  This can be "jh", and various others.
	GameMode string
	// The current level's name.
	CurrentLevel string
	// The seed used to generate the game.
	Seed uint32
}

func (f StoredSaveFile) AbsPath() string {
	return path.Join(getSaveScumDir(), f.StoredRelPath)
}

func (f StoredSaveFile) Delete(db *gorm.DB) {
	os.Remove(f.AbsPath())
	db.Delete(f)
}

func (f StoredSaveFile) Restore() error {
	return cp.Copy(
		f.AbsPath(),
		path.Join(guessGameDir(), f.OriginalBase),
	)
}
