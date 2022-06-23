package nvc

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

func ReadMagic(r io.Reader) error {
	header := [8]byte{}
	if _, err := io.ReadFull(r, header[:]); err != nil {
		return err
	}

	if string(header[:]) != MAGIC {
		return errors.New(".nvc signature not found")
	}

	return nil
}

func extractNVC(arcPath string, pathlist []string, outputDirectory string, extractUnknown bool) error {
	r, err := os.Open(arcPath)
	if err != nil {
		return err
	}

	if err :=  ReadMagic(r); err != nil {
		return err
	}

	hashToPath := map[uint64]string{}
	for _, p := range pathlist {
		hash := String2Hash(p)
		hashToPath[hash] = p
	}

	var tocEntries []TocEntry
	tocEntries, err = ReadToc(r)
	if err != nil {
		return err
	}

	entries := map[uint64]TocEntry{}
	for _, e := range(tocEntries) {
		entries[e.Hash] = e
	}

	extractedCount := 0

	for hash, nvcEntry := range entries {
		switch nvcEntry.Flags {
		case UNCOMPRESSED:
			// Do nothing
		case COMPRESSED:
			// TODO: Perform zlib decompression
		default:
			continue
		}

		_, exists := hashToPath[hash]
		if extractUnknown || exists {
			_, err = r.Seek(int64(nvcEntry.Offset), 0) // "0 means relative to the origin of the file"
			if err != nil {
				continue
			}

			data := make([]byte, nvcEntry.Length)
			_, err = io.ReadFull(r, data)
			if err != nil {
				continue
			}

			var path string

			if exists {
				path = hashToPath[hash]
			} else {
				var (
					dirName string
					ext     string
				)

				switch string(data[:4]) {
				case "\x89PNG":
					dirName = "data/unknown_png/"
					ext = ".png"
				case "nmf1":
					dirName = "data/unknown_nmd/"
					ext = ".nmd"
				case "OggS":
					dirName = "data/unknown_ogg/"
					ext = ".ogg"
				case "RIFF":
					dirName = "data/unknown_wav/"
					ext = ".wav"
				default:
					dirName = "data/unknown/"
					ext = ".unknown"
				}

				path = dirName + strconv.FormatUint(hash, 10) + ext
			}

			outputPath := outputDirectory + path

			outFile, err := os.Create(outputPath)
			if err != nil {
				continue
			}

			_, err = outFile.Write(data)
			if err != nil {
				continue
			}

			extractedCount++
		}
	}

	fmt.Printf("Extracted %d files\n", extractedCount)
	return nil
}
