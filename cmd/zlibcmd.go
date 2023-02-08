package cmd

import (
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/spf13/cobra"
)

var rawRE *regexp.Regexp
var zlibRE *regexp.Regexp

func init() {
	rawRE = regexp.MustCompile(`\.raw$`)
	zlibRE = regexp.MustCompile(`\.zlib$`)
}

func UnzlibCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unzlib FILE",
		Short: "Decompress a Zlib-compressed file.",
		Long: `Decompress a Zlib-compressed file.

If the input filename has a ".zlib" suffix, it is stripped and the result is
the output filename.  Otherwise the output filename is the input filename with
".raw" tacked on to the end.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				cmd.Usage()
				os.Exit(1)
			}
			inputFilename := args[0]
			outputFilename := rawRE.ReplaceAllString(inputFilename, "")
			if inputFilename == outputFilename {
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

func ZlibCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zlib FILE",
		Short: "Compress a file using Zlib.",
		Long: `Compress a file using Zlib.

If the FILE ends with .raw, this suffix is removed.
Otherwise a .zlib suffix is tacked to the end of the output filename.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				cmd.Usage()
				os.Exit(1)
			}
			inputFilename := args[0]
			outputFilename := rawRE.ReplaceAllString(inputFilename, "")
			if inputFilename == outputFilename {
				outputFilename = outputFilename + ".zlib"
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
