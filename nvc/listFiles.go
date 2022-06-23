package nvc

import (
	"io"
)

func ListEntries(r io.Reader) ([]TocEntry, error) {
	if magicErr := ReadMagic(r); magicErr != nil {
		return nil, magicErr
	}

	return ReadToc(r)
}
