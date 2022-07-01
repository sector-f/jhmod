package savedb

import (
	"database/sql"
	"errors"
	"io"

	"github.com/sector-f/jhmod/savefile"
)

type DB struct {
	db *sql.DB
}

func New(filename string) (*DB, error) {
	return nil, errors.New("unimplemented function")
}

func (d *DB) init() error {
	return errors.New("unimplemented function")
}

func (d *DB) Insert(s savefile.SaveData, r io.ReadSeeker) error {
	return errors.New("unimplemented function")
}

func (d *DB) List() ([]savefile.SaveData, error) {
	return []savefile.SaveData{}, errors.New("unimplemented function")
}

func (d *DB) Delete(id int) error {
	return errors.New("unimplemented function")
}

type saveModel struct {
	id         int // Should this be a UUID? Hash of file data?
	playerName string
	gameType   string
	curLevel   string
	seed       uint32
	file       []byte
}
