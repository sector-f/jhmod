package cmd

import (
	"fmt"
	"os"

	"github.com/sector-f/jh_extract/nvc"
)

func extractNVC(arcPath string, pathlist []string, outputDirectory string, extractUnknown bool) error {
	r, err := os.Open(arcPath)
	if err != nil {
		return err
	}

	hashToPath := map[nvc.Hash]string{}
	for _, p := range pathlist {
		hash := nvc.String2Hash(p)
		hashToPath[hash] = p
	}

	tocEntries, err := nvc.ReadToc(r)
	if err != nil {
		return err
	}

	entries := map[nvc.Hash]nvc.TocEntry{}
	for _, e := range tocEntries {
		entries[e.Hash] = e
	}

	extractedCount := 0

	for hash, nvcEntry := range entries {
		_, exists := hashToPath[hash]
		if extractUnknown || exists {
			data, err := nvcEntry.Data(r)
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

				path = fmt.Sprintf("%s%v%s", dirName, hash, ext)
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
