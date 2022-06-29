package savefile

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"io"
	"io/ioutil"

	"github.com/sector-f/jhmod/nvc"
)

const (
	// Magic header for Jupiter Hell save save files.
	magic = "\xde\xc0\xad\xde"
)

type savedata struct {
	// Player name
	PlayerName string
	// Game mode.  This can be "jh", and various others.
	GameMode string
	// The current level's name.
	CurrentLevel string
	// The seed used to generate the game.
	Seed uint32
}

// Something was reading for data but got an EOF.
var UnexpectedEndOfFileErr error = errors.New("unexpected eof")

// General parse error.  TODO make more descriptive errors.
var SaveParseErr error = errors.New("general parse error")

// Read in a string from a buffer of size suffixed with a length byte.
func readString(r io.Reader, size int) (string, error) {
	buf := make([]byte, size)
	nRead, err := r.Read(buf[:])
	if err != nil {
		return "", err
	}
	if nRead != size {
		return "", UnexpectedEndOfFileErr
	}
	recordedSize := buf[size-1]
	if int(recordedSize) > size {
		return "", SaveParseErr
	}
	stringBytes := buf[:recordedSize]

	for _, b := range stringBytes {
		if b == 0 {
			return "", SaveParseErr
		}
	}
	// TODO verify valid text.

	return string(stringBytes), nil
}

// Advance the reader r by size bytes.
func discard(r io.Reader, size int) error {
	buf := make([]byte, size)
	_, err := io.ReadFull(r, buf[:])
	return err
}

// Attempt to read in a save file.
//
// If the error value is non-nil, the return value does not contain valid data.
func Parse(r io.Reader) (savedata, error) {
	z, zErr := zlib.NewReader(r)
	if zErr != nil {
		return savedata{}, zErr
	}
	uncompressed, uErr := ioutil.ReadAll(z)
	if uErr != nil {
		return savedata{}, uErr
	}
	b := bytes.NewReader(uncompressed)

	header := [len(magic)]byte{}
	if _, err := io.ReadFull(b, header[:]); err != nil {
		return savedata{}, err
	}
	if string(header[:]) != magic {
		return savedata{}, nvc.ErrNoMagicFound
	}

	// Not sure what these bytes are, sometimes they're:
	//  a900 0000
	//  a200 0000
	if err := discard(b, 4); err != nil {
		return savedata{}, err
	}

	gameMode, gameModeErr := readString(b, 32)
	if gameModeErr != nil {
		return savedata{}, gameModeErr
	}

	name, nameErr := readString(b, 32)
	if nameErr != nil {
		return savedata{}, nameErr
	}

	cur, curErr := readString(b, 64)
	if curErr != nil {
		return savedata{}, curErr
	}

	// Not sure what these 32 bytes are...
	if err := discard(b, 32); err != nil {
		return savedata{}, err
	}

	seed := uint32(0)
	if err := binary.Read(b, binary.LittleEndian, &seed); err != nil {
		return savedata{}, err
	}

	return savedata{
		GameMode:     gameMode,
		PlayerName:   name,
		CurrentLevel: cur,
		Seed:         seed,
	}, nil
}
