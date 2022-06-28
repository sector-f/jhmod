package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/sector-f/jh_extract/nvc"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func createCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create ARCHIVE [FILE]...",
		Short: "Create a .nvc archive",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			verbose, _ := cmd.PersistentFlags().GetBool("verbose")

			shouldCompress := false
			cmd.Flags().Visit(func(f *pflag.Flag) {
				if f.Name == "compress" {
					shouldCompress = true
				}
			})
			compressLevel, _ := cmd.PersistentFlags().GetInt("compress")
			if shouldCompress {
				if compressLevel < 0 || compressLevel > 9 {
					return errors.New("Compression level must be between 0-9")
				}
			}

			arcFilename := args[0]
			arcFile, err := os.Create(arcFilename)
			if err != nil {
				return err
			}
			defer arcFile.Close()

			fileNames := args[1:]
			writer, err := nvc.NewWriter(arcFile, uint32(len(fileNames)))
			if err != nil {
				return err
			}

			for _, fName := range fileNames {
				if verbose {
					fmt.Println(fName)
				}

				file, err := os.Open(fName)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error opening %s: %v\n", fName, err)
					continue
				}
				defer file.Close()

				hashedName := nvc.String2Hash(fName)

				if shouldCompress {
					_, err = writer.CreateCompressed(file, hashedName, compressLevel)
				} else {
					_, err = writer.Create(file, hashedName)
				}
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error adding %s: %v\n", fName, err)
					continue
				}
			}

			return writer.Finalize()
		},
	}

	cmd.PersistentFlags().BoolP("verbose", "v", false, "Print the names of files to standard output")
	cmd.PersistentFlags().IntP("compress", "c", 0, "Compression level 0-9 (where 0 is no compression, 1 is best speed, and 9 is best compression)")

	return cmd
}
