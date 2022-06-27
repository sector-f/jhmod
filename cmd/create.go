package cmd

import (
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

			isFlagPassed := func(name string) bool {
				found := false
				cmd.PersistentFlags().Visit(func(f *pflag.Flag) {
					if f.Name == name {
						found = true
					}
				})
				return found
			}

			shouldCompress := isFlagPassed("compress")
			compressLevel, _ := cmd.PersistentFlags().GetInt("verbose")

			arcFilename := args[0]
			arcFile, err := os.Create(arcFilename)
			if err != nil {
				return err
			}

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
	cmd.PersistentFlags().IntP("compress", "c", 0, "Compression level")

	return cmd
}
