package gui

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	path "path/filepath"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/sector-f/jhmod/savefile"
	"github.com/u-root/u-root/pkg/cp"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm" // Base ORM
)

func sha256File(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return make([]byte, 0), err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return make([]byte, 0), err
	}
	return h.Sum(nil), nil
}

func Run() {
	a := app.New()
	w := a.NewWindow("jhmod manager")
	lbl := widget.NewLabel("Last Save Will Show here")
	var last *StoredSaveFile = nil
	restoreLast := widget.NewButton("Restore Last", func() {
		if last != nil {
			last.Restore()
		}
	})
	updateLast := func(f *StoredSaveFile) {
		last = f
		if f == nil {
			lbl.SetText("N/A")
			restoreLast.SetText("Restore N/A")
			restoreLast.Disable()
		} else {
			restoreLast.SetText(fmt.Sprintf("Restore %s", f.OriginalBase))
			restoreLast.Enable()
			lbl.SetText(fmt.Sprintf("[%s] %s",
				f.CreatedAt.Format("2006-01-02 15:04:06"),
				f.OriginalBase))
		}
	}
	content := container.New(
		layout.NewVBoxLayout(),
		widget.NewLabel("jhmod manager"),
		lbl,
		restoreLast,
	)

	mkdirErr := os.Mkdir(getSaveScumDir(), 0750)
	if mkdirErr != nil && !os.IsExist(mkdirErr) {
		fmt.Fprintf(os.Stderr, "Cannot make savescumdir %v\n", mkdirErr)
		return
	}

	dbPath := path.Join(
		getSaveScumDir(),
		"db.sqlite3",
	)
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Unable to open db \"%s\" %v", dbPath, err))
	}
	if err = db.AutoMigrate(&StoredSaveFile{}); err != nil {
		panic(fmt.Sprintf("Could not migrate db.\n"))
	}

	db.Last(&last)
	updateLast(last)

	go watchForNewSave(func(p string) {
		sd, err := savefile.ParseFile(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not parse save \"%s\".  Aborting.\n", p)
			return
		}
		digest, shaErr := sha256File(p)
		if shaErr != nil {
			fmt.Fprintf(os.Stderr, "Could not SHA256 \"%s\".  Aborting.\n", p)
			return
		}

		digestHex := hex.EncodeToString(digest)
		relPath := digestHex
		var existing *StoredSaveFile = nil
		if db.Where("sha256_hex = ?", digestHex).Find(&existing).RowsAffected > 0 && existing != nil {
			fmt.Fprintf(os.Stderr, "Saw save we already have in DB, skipping (%s).\n", digestHex)
			return
		}

		storedSaveFile := &StoredSaveFile{
			OriginalBase:  path.Base(p),
			StoredRelPath: relPath,
			Sha256Hex:     digestHex,

			PlayerName:   sd.PlayerName,
			GameMode:     sd.GameMode,
			CurrentLevel: sd.CurrentLevel,
			Seed:         sd.Seed,
		}

		destAbs := path.Join(getSaveScumDir(), relPath)
		if cpErr := cp.Copy(p, destAbs); cpErr != nil {
			fmt.Fprintf(os.Stderr, "Failed to copy file \"%s\" to \"%s\" (%v)\n", p, destAbs, cpErr)
			return
		} else {
			fmt.Printf("Copied \"%s\" to \"%s\" n saved to db.\n", p, destAbs)
		}
		db.Create(storedSaveFile)
		updateLast(storedSaveFile)
	})

	w.SetContent(content)
	w.ShowAndRun()
}
