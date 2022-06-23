package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"os"
)

// TODO: maybe pick a better name for this type
type compressionType uint32

const (
	uncompressed compressionType = 0
	compressed   compressionType = 1
	encrypted    compressionType = 3

	nvcMagicBytes = "nvc1d\x00\x00\x00"
)

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

	return errors.New("unimplemented")
}

func main() {
	fmt.Println("vim-go")
}
