package nvc

import (
	"fmt"
	"io"
	"errors"
	"encoding/binary"
)

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
	if magicErr := readMagic(r) ; magicErr != nil {
		return nil, magicErr
	}
	count, countErr := readCount(r)
	if countErr != nil {
		return nil, countErr
	}

	entries := make([]TocEntry, count)

	for i := range(entries) {
		entry, eErr := readEntry(r)
		if eErr != nil {
			return nil, eErr
		}
		entries[i] = entry
	}

	return entries, nil

}

// Read the magic bytes at the beginning of a .nvc file.
func readMagic(r io.Reader) error {
	header := [8]byte{}
	if _, err := io.ReadFull(r, header[:]); err != nil {
		return err
	}

	if string(header[:]) != MAGIC {
		return errors.New(".nvc signature not found")
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
