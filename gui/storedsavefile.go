package gui

import (
	"fmt"
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

	Restores []RestoreRecord
}

type RestoreRecord struct {
	gorm.Model
	Id               int `gorm:"primary_key"`
	StoredSaveFileID int `gorm:"foreign_key:StoredSaveFile"`
}

func (f StoredSaveFile) String() string {
	return fmt.Sprintf("% 3d [%s] %s %s %s (%d)",
		f.Id,
		f.CreatedAt.Format("2006-01-02 15:04:06"),
		f.GameMode,
		f.PlayerName,
		f.CurrentLevel,
		f.Seed,
	)
}

func (f StoredSaveFile) AbsPath() string {
	return path.Join(getSaveScumDir(), f.StoredRelPath)
}

func (f StoredSaveFile) Delete(db *gorm.DB) {
	os.Remove(f.AbsPath())
	db.Delete(f)
}

func (f StoredSaveFile) Restore() (*RestoreRecord, error) {
	err := cp.Copy(
		f.AbsPath(),
		path.Join(guessGameDir(), f.OriginalBase),
	)
	if err != nil {
		return nil, err
	}
	return &RestoreRecord{StoredSaveFileID: f.Id}, nil
}
