// Package nvc implements a parser for nvc archive files.
//
// The format of an archive is described below. All header values are stored on disk as little endian.
//
//   1. NVC magic bytes (8 bytes)
//   2. Number of entries in the archive's table of contents (32-bit little endian unsigned integer)
//   3. Table of Contents entries (192 * n bytes, where n is the previously-specified number of entries)
package nvc

import (
	"compress/zlib"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
)

// EntryFlags describes how a file is stored on disk
type EntryFlags uint32

const (
	// EntryFlagNoCompression indicates that the file is stored in the NVC archive uncompressed
	EntryFlagNoCompression EntryFlags = 0
	// EntryFlagZlibCompression indicates that the file is stored in the NVC archive with zlib compression
	EntryFlagZlibCompression EntryFlags = 1
	// EntryFlagEncrypted indicates that the file is encrypted
	EntryFlagEncrypted EntryFlags = 3

	// NVC file type magic bytes
	magic = "nvc1d\x00\x00\x00"
)

var ErrNoMagicFound error = errors.New("nvc magic bytes not found")

// TocEntry is a file entry in the table of contents.
type TocEntry struct {
	Hash      Hash       // 64-bit FNV-1a hash of the file's path on disk
	Offset    uint32     // File's offset (in bytes) from the start of the nvc file
	RawLength uint32     // Length (in bytes) of the file after it has been extracted
	Length    uint32     // Length (in bytes) of file as it is stored in the nvc file
	Flags     EntryFlags // Indicates whether file is compressed or encrypted (TODO: determine if this is a bitmask)
}

// ReadToc parses r as an NVC archive and returns its table of contents.
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

// Data returns from r the file that t is an entry for.
func (t TocEntry) Data(r io.ReadSeeker) ([]byte, error) {
	_, err := r.Seek(int64(t.Offset), 0)
	if err != nil {
		return nil, err
	}

	var reader io.Reader = r

	switch t.Flags {
	case EntryFlagNoCompression:
		// Do nothing
	case EntryFlagZlibCompression:
		reader, err = zlib.NewReader(r)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unsupported")
	}

	data := make([]byte, t.RawLength)
	_, err = io.ReadFull(reader, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (e TocEntry) String() string {
	return fmt.Sprintf("%v offset=%v %vB (%vB on disk) flags=%v",
		e.Hash,
		e.Offset,
		e.RawLength,
		e.Length,
		e.Flags)
}

// readMagic reads the magic bytes at the beginning of a .nvc file.
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

// readEntry reads a single ToC entry from r.
func readEntry(r io.Reader) (TocEntry, error) {
	var e TocEntry
	if err := binary.Read(r, binary.LittleEndian, &e); err != nil {
		return e, err
	}
	return e, nil
}

// readCount reads the ToC entry count from r.
func readCount(r io.Reader) (uint32, error) {
	var count uint32
	if err := binary.Read(r, binary.LittleEndian, &count); err != nil {
		return 0, err
	}
	return count, nil
}

// Hash is a 64-bit FNV-1a hash
type Hash uint64

// String2Hash returns the 64-bit FNV-1a hash of s
func String2Hash(s string) Hash {
	hash := fnv.New64a()
	hash.Write([]byte(s))
	return Hash(hash.Sum64())
}

func (h Hash) String() string {
	return fmt.Sprintf("%016x", uint64(h))
}
