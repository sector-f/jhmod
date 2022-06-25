package nvc

import (
	"compress/zlib"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type EntryFlags uint32

const (
	Uncompressed EntryFlags = 0
	Compressed   EntryFlags = 1
	Encrypted    EntryFlags = 3

	magic = "nvc1d\x00\x00\x00"
)

var ErrNoMagicFound error = errors.New("nvc magic bytes not found")

type TocEntry struct {
	Hash      Hash
	Offset    uint32
	RawLength uint32 // Uncompressed length
	Length    uint32 // Actual length in file
	Flags     EntryFlags
}

// Reads the Magic Header and Table of Contents from r.
// The ToC format is as follows:
//   1. Magic bytes header (technically not part of the ToC but it comes first)
//   2. 32 LE unsigned integer -> number of entries in ToC.
//   3. This repeated as many times as the last value indicates:
//     a. FNV1a hash (64-bit LE unsigned integer)
//     b. Offset in the NVC archive (64-bit LE unsigned integer)
//     c. Uncompressed length of the Entry (64-bit LE unsigned integer)
//     d. Actual length of the Entry (64-bit LE unsigned integer)
//     e. Entry flags (64-bit LE unsigned integer)
func ReadToc(r io.Reader) ([]TocEntry, error) {
	if magicErr := readMagic(r); magicErr != nil {
		return nil, magicErr
	}
	count, countErr := readCount(r)
	if countErr != nil {
		return nil, countErr
	}

	entries := make([]TocEntry, count)

	for i := range entries {
		entry, eErr := readEntry(r)
		if eErr != nil {
			return nil, eErr
		}
		entries[i] = entry
	}

	return entries, nil

}

func (t TocEntry) Data(r io.ReadSeeker) ([]byte, error) {
	_, err := r.Seek(int64(t.Offset), 0)
	if err != nil {
		return []byte{}, err
	}

	var reader io.Reader = r

	switch t.Flags {
	case Uncompressed:
		// Do nothing
	case Compressed:
		reader, err = zlib.NewReader(r)
		if err != nil {
			return []byte{}, err
		}
	default:
		return []byte{}, errors.New("unsupported")
	}

	data := make([]byte, t.RawLength)
	_, err = io.ReadFull(reader, data)
	if err != nil {
		return []byte{}, err
	}

	return data, nil
}

// Read the magic bytes at the beginning of a .nvc file.
func readMagic(r io.Reader) error {
	header := [8]byte{}
	if _, err := io.ReadFull(r, header[:]); err != nil {
		return err
	}

	if string(header[:]) != magic {
		return ErrNoMagicFound
	}

	return nil
}

// Read a single entry from the ToC.
func readEntry(r io.Reader) (TocEntry, error) {
	var e TocEntry
	if err := binary.Read(r, binary.LittleEndian, &e); err != nil {
		return e, err
	}
	return e, nil
}

// Read the ToC entry count.
func readCount(r io.Reader) (uint32, error) {
	var count uint32
	if err := binary.Read(r, binary.LittleEndian, &count); err != nil {
		return 0, err
	}
	return count, nil
}

func (e TocEntry) String() string {
	return fmt.Sprintf("%v offset=%v %vB (%vB on disk) flags=%v",
		e.Hash,
		e.Offset,
		e.RawLength,
		e.Length,
		e.Flags)
}
