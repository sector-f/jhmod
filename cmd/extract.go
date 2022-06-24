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
			verbose, _ := cmd.PersistentFlags().GetBool("verbose")

			pathFile, err := os.Open(pathFilename)
			if err != nil {
				return err
			}
			defer pathFile.Close()

			pathlist := []string{}
			scanner := bufio.NewScanner(pathFile)
			for scanner.Scan() {
				pathlist = append(pathlist, scanner.Text())
			}

			return extractNVC(arcFilename, pathlist, outputDir, extractUnknown, verbose)
		},
	}

	cmd.PersistentFlags().StringP("file", "f", "", "Path to NVC file")
	cmd.PersistentFlags().StringP("pathlist", "p", "", "Path to pathlist file")
	cmd.PersistentFlags().StringP("output", "o", "", "Output directory")
	cmd.PersistentFlags().BoolP("unknown", "u", false, "Additionally files which are not named in the pathlist file")
	cmd.PersistentFlags().BoolP("verbose", "v", false, "Print the names of extracted files to standard output")

	return cmd
}

func extractNVC(arcPath string, pathlist []string, outputDirectory string, extractUnknown bool, verbose bool) error {
	arcFile, err := os.Open(arcPath)
	if err != nil {
		return err
	}
	defer arcFile.Close()

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
				type filetype struct {
					dirName string
					ext     string
				}

				filetypes := map[string]filetype{
					"\x89PNG":          {"unknown_png", ".png"},
					"nmf1":             {"unknown_nmd", ".nmd"},
					"OggS":             {"unknown_ogg", ".ogg"},
					"RIFF":             {"unknown_wav", ".wav"},
					"\x03\x02\x23\x07": {"unknown_spirv", ".spirv"},
				}

				// Default values that get overwritten if possible
				var (
					dirName string = "unknown"
					ext     string = ".unknown"
				)

				magicBytes := string(data[:4])
				if ftype, exists := filetypes[magicBytes]; exists {
					dirName = ftype.dirName
					ext = ftype.ext
				}

				path = filepath.Join("data", dirName, hash.String()+ext)
			}

			outputPath := filepath.Join(outputDirectory, path)

			if verbose {
				fmt.Printf("Extracting %s to %s\n", hash.String(), outputPath)
			}

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
			defer outFile.Close()

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
