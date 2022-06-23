package nvc

import (
	"io"
	"encoding/binary"
)

// Reads the Table of Contents from r.
// The ToC format is as follows:
//   1. 32 LE unsigned integer -> number of entries in ToC.
//   2. This repeated as many times as the last value indicates:
//     a. FNV1a hash (64-bit LE unsigned integer)
//     b. Offset in the NVC archive (64-bit LE unsigned integer)
//     c. Uncompressed length of the Entry (64-bit LE unsigned integer)
//     d. Actual length of the Entry (64-bit LE unsigned integer)
//     e. Entry flags (64-bit LE unsigned integer)
func ReadToc(r io.Reader) ([]TocEntry, error) {
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
