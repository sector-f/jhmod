package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"

	"github.com/spf13/cobra"
)

func init() {
	pathlistCmd.AddCommand(pathlistScanCmd())
}

var pathlistCmd = &cobra.Command{
	Use:   "pathlist",
	Short: "Work with pathlist files",
}

func pathlistScanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scan FILE",
		Short: "Scan a core dump for interesting paths",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			verbose, _ := cmd.PersistentFlags().GetBool("verbose")
			regex := *regexp.MustCompile("data/[A-Za-z0-9_./]+\\.[a-zA-Z0-9]+")

			reader, openErr := os.Open(args[0])
			if openErr != nil {
				fmt.Fprintln(os.Stderr, openErr)
				os.Exit(1)
			}
			defer reader.Close()
			r := bufio.NewReader(reader)

			matches := make(map[string]struct{})
			offset := int64(0)
			for { // Inspired from https://go.dev/play/p/aPrAW7XGHi

				m := regex.FindReaderIndex(r)
				if m == nil {
					break
				}

				offset += int64(m[0])  // Add the offset.
				reader.Seek(offset, 0) // Seek to the match.
				r.Reset(reader)        // Flush buffer so buffered reader is looking where we want it to.

				b := make([]byte, m[1]-m[0])
				io.ReadFull(r, b)
				s := string(b[:])

				if verbose {
					fmt.Fprintf(os.Stderr, "offset=%v s=%v\n", offset, s)
				}

				matches[s] = struct{}{}      // Set membership
				offset += int64(m[1] - m[0]) // Keep record of where we think the reader is looking next.
			}

			if verbose {
				fmt.Fprintf(os.Stderr, "\nFound %v matches.\n", len(matches))
			}

			ary := []string{}
			for path, _ := range matches {
				ary = append(ary, path)
			}
			sort.Strings(ary)

			for _, path := range ary {
				fmt.Println(path)
			}
		},
	}
	cmd.PersistentFlags().BoolP("verbose", "v", false, "Print scan data in realtime.  Print summary at end.")

	return cmd
}
