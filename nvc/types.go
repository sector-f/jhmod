package nvc

import (
	"compress/zlib"
	"errors"
	"io"
)

type EntryFlags uint32
type Hash uint64

const (
	Uncompressed EntryFlags = 0
	Compressed   EntryFlags = 1
	Encrypted    EntryFlags = 3

	magic = "nvc1d\x00\x00\x00"
)

type TocEntry struct {
	Hash      Hash
	Offset    uint32
	RawLength uint32 // Uncompressed length
	Length    uint32 // Actual length in file
	Flags     EntryFlags
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
