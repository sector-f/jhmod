package savefile

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"io"
	"io/ioutil"
	"unicode"
)

const (
	// Magic header for Jupiter Hell save save files.
	magic = "\xde\xc0\xad\xde"
)

type SaveData struct {
	// Player name
	PlayerName string
	// Game mode.  This can be "jh", and various others.
	GameMode string
	// The current level's name.
	CurrentLevel string
	// The seed used to generate the game.
	Seed uint32
}

// No magic found on save file.
var ErrNoMagicFound error = errors.New("no magic found")

// General parse error.  TODO make more descriptive errors.
var SaveParseErr error = errors.New("general parse error")

// Return the last byte of buf or panic() when not possible.
func lastByte(buf []byte) byte {
	if len(buf) == 0 {
		panic("lastByte with zero length slice")
	}
	return buf[len(buf)-1]
}

// Returns true if byte looks like ascii, false otherwise.
func isAscii(b byte) bool {
	return b > 0 && b <= unicode.MaxASCII
}

// readString reads size bytes from r, interprets the final byte read as an unsigned integer `n`, and returns a string consisting of the first n bytes that were read from r
func readString(r io.Reader, size uint) (string, error) {
	buf := make([]byte, size)
	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return "", err
	}
	recordedSize := lastByte(buf)
	if uint(recordedSize) > size {
		return "", SaveParseErr
	}
	stringBytes := buf[:recordedSize]

	for _, b := range stringBytes {
		if !isAscii(b) {
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
func Parse(r io.Reader) (SaveData, error) {
	z, zErr := zlib.NewReader(r)
	if zErr != nil {
		return SaveData{}, zErr
	}
	uncompressed, uErr := ioutil.ReadAll(z)
	if uErr != nil {
		return SaveData{}, uErr
	}
	b := bytes.NewReader(uncompressed)

	header := [len(magic)]byte{}
	if _, err := io.ReadFull(b, header[:]); err != nil {
		return SaveData{}, err
	}
	if string(header[:]) != magic {
		return SaveData{}, ErrNoMagicFound
	}

	// Not sure what these bytes are, sometimes they're:
	//  a900 0000
	//  a200 0000
	if err := discard(b, 4); err != nil {
		return SaveData{}, err
	}

	gameMode, gameModeErr := readString(b, 32)
	if gameModeErr != nil {
		return SaveData{}, gameModeErr
	}

	name, nameErr := readString(b, 32)
	if nameErr != nil {
		return SaveData{}, nameErr
	}

	cur, curErr := readString(b, 64)
	if curErr != nil {
		return SaveData{}, curErr
	}

	// Not sure what these 32 bytes are...
	if err := discard(b, 32); err != nil {
		return SaveData{}, err
	}

	seed := uint32(0)
	if err := binary.Read(b, binary.LittleEndian, &seed); err != nil {
		return SaveData{}, err
	}

	return SaveData{
		GameMode:     gameMode,
		PlayerName:   name,
		CurrentLevel: cur,
		Seed:         seed,
	}, nil
}
