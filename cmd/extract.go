package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sector-f/jh_extract/nvc"
	"github.com/spf13/cobra"
)

func extractCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extract",
		Short: "Extract files from a .nvc file",
		RunE: func(cmd *cobra.Command, args []string) error {
			arcFilename, _ := cmd.PersistentFlags().GetString("file")
			pathFilename, _ := cmd.PersistentFlags().GetString("pathlist")
			outputDir, _ := cmd.PersistentFlags().GetString("output")
			extractUnknown, _ := cmd.PersistentFlags().GetBool("unknown")

			pathFile, err := os.Open(pathFilename)
			if err != nil {
				return err
			}

			pathlist := []string{}
			scanner := bufio.NewScanner(pathFile)
			for scanner.Scan() {
				pathlist = append(pathlist, scanner.Text())
			}

			return extractNVC(arcFilename, pathlist, outputDir, extractUnknown)
		},
	}

	cmd.PersistentFlags().StringP("file", "f", "", "Path to NVC file")
	cmd.PersistentFlags().StringP("pathlist", "p", "", "Path to pathlist file")
	cmd.PersistentFlags().StringP("output", "o", "", "Output directory")
	cmd.PersistentFlags().BoolP("unknown", "u", false, "Additionally files which are not named in the pathlist file")

	return cmd
}

func extractNVC(arcPath string, pathlist []string, outputDirectory string, extractUnknown bool) error {
	arcFile, err := os.Open(arcPath)
	if err != nil {
		return err
	}

	hashToPath := map[nvc.Hash]string{}
	for _, p := range pathlist {
		hash := nvc.String2Hash(p)
		hashToPath[hash] = p
	}

	tocEntries, err := nvc.ReadToc(arcFile)
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
			data, err := nvcEntry.Data(arcFile)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}

			var path string // File path/name as stored (hashed) in ToC

			if exists { // If this entry's name was in pathlist, then we use its name...
				path = hashToPath[hash]
			} else { // ...And if it wasn't, then we make up a name using the hash
				var (
					dirName string
					ext     string
				)

				switch string(data[:4]) {
				case "\x89PNG":
					dirName = filepath.Join("data", "unknown_png")
					ext = ".png"
				case "nmf1":
					dirName = filepath.Join("data", "unknown_nmd")
					ext = ".nmd"
				case "OggS":
					dirName = filepath.Join("data", "unknown_ogg")
					ext = ".ogg"
				case "RIFF":
					dirName = filepath.Join("data", "unknown_wav")
					ext = ".wav"
				default:
					dirName = filepath.Join("data", "unknown")
					ext = ".unknown"
				}

				path = filepath.Join(dirName, hash.String()+ext)
			}

			outputPath := filepath.Join(outputDirectory, path)

			err = os.MkdirAll(filepath.Dir(outputPath), 0755)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}

			outFile, err := os.Create(outputPath)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}

			_, err = outFile.Write(data)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}

			extractedCount++
		}
	}

	fmt.Printf("Extracted %d files\n", extractedCount)
	return nil
}
