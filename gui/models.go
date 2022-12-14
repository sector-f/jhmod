package gui

import (
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
