package cmd

import (
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func unzlibCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unzlib FILE",
		Short: "Decompress a Zlib-compressed file.",
		Long: `Decompress a Zlib-compressed file.

If the input filename has a ".zlib" suffix, it is stripped and the result is
the output filename.  Otherwise the output filename is the input filename with
".raw" tacked on to the end.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inputFilename := args[0]
			var outputFilename string

			if strings.HasSuffix(inputFilename, ".zlib") {
				outputFilename = strings.TrimSuffix(inputFilename, ".zlib")
			} else {
				outputFilename = inputFilename + ".raw"
			}

			fmt.Fprintf(os.Stderr, "Writing to %s\n", outputFilename)

			fi, err1 := os.Open(inputFilename)
			if err1 != nil {
				return err1
			}
			defer fi.Close()

			fo, err2 := os.OpenFile(outputFilename, os.O_CREATE|os.O_WRONLY, 0755)
			if err2 != nil {
				return err2
			}
			defer fo.Close()

			reader, err3 := zlib.NewReader(fi)
			defer reader.Close()
			if err3 != nil {
				panic(err3)
			}

			_, err4 := io.Copy(fo, reader)
			if err4 != nil {
				panic(err4)
			}

			return nil
		},
	}
	return cmd
}

func zlibCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zlib FILE",
		Short: "Compress a file using Zlib.",
		Long: `Compress a file using Zlib.

If the FILE ends with .raw, this suffix is removed.
Otherwise a .zlib suffix is tacked to the end of the output filename.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inputFilename := args[0]
			var outputFilename string

			if strings.HasSuffix(inputFilename, ".raw") {
				outputFilename = strings.TrimSuffix(inputFilename, ".raw")
			} else {
				outputFilename = inputFilename + ".zlib"
			}

			fmt.Fprintf(os.Stderr, "Writing to %s\n", outputFilename)

			fi, err1 := os.Open(inputFilename)
			if err1 != nil {
				return err1
			}
			defer fi.Close()

			fo, err2 := os.OpenFile(outputFilename, os.O_CREATE|os.O_WRONLY, 0755)
			if err2 != nil {
				return err2
			}
			defer fo.Close()

			writer := zlib.NewWriter(fo)
			defer writer.Close()
			_, err4 := io.Copy(writer, fi)
			if err4 != nil {
				panic(err4)
			}

			return nil
		},
	}
	return cmd
}
