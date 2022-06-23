package nvc

import (
	"encoding/binary"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"os"
	"strconv"
)

// TODO: maybe pick a better name for this type
type compressionType uint32

const (
	uncompressed compressionType = 0
	compressed   compressionType = 1
	encrypted    compressionType = 3

	nvcMagicBytes = "nvc1d\x00\x00\x00"
)

type nvcEntry struct {
	Hash        uint64
	Offset      uint32
	RawLength   uint32 // Uncompressed length
	Length      uint32 // Actual length in file
	Compression compressionType
}

func Hash2String(hash uint64) string {
	return fmt.Sprintf("%016x", hash)
}

func String2Hash(s string) uint64 {
	hash := fnv.New64a()
	hash.Write([]byte(s))
	return hash.Sum64()
}

func readNVCMagic(r io.Reader) (bool, error) {
	header := [8]byte{}
	if _, err := io.ReadFull(r, header[:]); err != nil {
		return false, err
	}

	if string(header[:]) != nvcMagicBytes {
		return false, errors.New(".nvc signature not found")
	}

	return true, nil
}

func readNvcEntry(r io.Reader) (nvcEntry, error) {
	var e nvcEntry
	if err := binary.Read(r, binary.LittleEndian, &e); err != nil {
		return e, err
	}
	return e, nil
}

func readCount(r io.Reader) (uint32, error) {
	var count uint32
	if err := binary.Read(r, binary.LittleEndian, &count); err != nil {
		return 0, err
	}
	return count, nil
}

func extractNVC(arcPath string, pathlist []string, outputDirectory string, extractUnknown bool) error {
	arcFile, err := os.Open(arcPath)
	if err != nil {
		return err
	}

	if ok, err :=  readNVCMagic(arcFile); !ok {
		return err
	}

	hashToPath := map[uint64]string{}
	for _, p := range pathlist {
		hash := String2Hash(p)
		hashToPath[hash] = p
	}

	var numEntries uint32
	numEntries, err = readCount(arcFile)
	if err != nil {
		return nil
	}

	entries := map[uint64]nvcEntry{}

	for i := 0; i < int(numEntries); i++ {
		var e nvcEntry
		binary.Read(arcFile, binary.LittleEndian, &e)
		entries[e.Hash] = e
	}

	extractedCount := 0

	for hash, nvcEntry := range entries {
		switch nvcEntry.Compression {
		case uncompressed:
			// Do nothing
		case compressed:
			// TODO: Perform zlib decompression
		default:
			continue
		}

		_, exists := hashToPath[hash]
		if extractUnknown || exists {
			_, err = arcFile.Seek(int64(nvcEntry.Offset), 0) // "0 means relative to the origin of the file"
			if err != nil {
				continue
			}

			data := make([]byte, nvcEntry.Length)
			_, err = io.ReadFull(arcFile, data)
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
