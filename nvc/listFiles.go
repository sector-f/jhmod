package nvc

import (
	"io"
)

func ListEntries(r io.Reader) ([]nvcEntry, error) {
	if ok, magicErr := readNVCMagic(r); !ok {
		return nil, magicErr
	}

	count, countErr := readCount(r)
	if countErr != nil {
		return nil, countErr
	}

	entries := make([]nvcEntry, count)

	for i := range(entries) {
		entry, eErr := readNvcEntry(r)
		if eErr != nil {
			return nil, eErr
		}
		entries[i] = entry
	}

	return entries, nil
}
