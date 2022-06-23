package main

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

type entry struct {
	Hash        uint64
	Offset      uint32
	_           uint32 // Uncompressed length
	Length      uint32 // Actual length in file
	Compression compressionType
}

func extractNVC(arcPath string, pathlist []string, outputDirectory string, extractUnknown bool) error {
	arcFile, err := os.Open(arcPath)
	if err != nil {
		return err
	}

	header := [8]byte{}
	_, err = io.ReadFull(arcFile, header[:])
	if err != nil {
		return err
	}

	if string(header[:]) != nvcMagicBytes {
		return errors.New(".nvc signature not found")
	}

	hashToPath := map[uint64]string{}
	for _, p := range pathlist {
		hash := fnv.New64a()
		hash.Write([]byte(p))
		hashToPath[hash.Sum64()] = p
	}

	var numEntries uint32
	err = binary.Read(arcFile, binary.LittleEndian, &numEntries)
	if err != nil {
		return nil
	}

	entries := map[uint64]entry{}

	for i := 0; i < int(numEntries); i++ {
		var e entry
		binary.Read(arcFile, binary.LittleEndian, &e)
		entries[e.Hash] = e
	}

	extractedCount := 0

	for hash, entry := range entries {
		switch entry.Compression {
		case uncompressed:
			// Do nothing
		case compressed:
			// TODO: Perform zlib decompression
		default:
			continue
		}

		_, exists := hashToPath[hash]
		if extractUnknown || exists {
			_, err = arcFile.Seek(int64(entry.Offset), 0) // "0 means relative to the origin of the file"
			if err != nil {
				continue
			}

			data := make([]byte, entry.Length)
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

func main() {
	fmt.Println("vim-go")
}
