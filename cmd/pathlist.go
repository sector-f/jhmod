package cmd

import 	(
	"fmt"
	"os"
	"regexp"
	"bufio"
	"io"
	"sort"

	"github.com/spf13/cobra"
)

const (
	regex = "data/[A-Za-z0-9_./]+\\.[a-zA-Z0-9]+"
)

var pattern regexp.Regexp

func init() {
	pathlistCmd.AddCommand(pathlistScanCmd())
	pattern = *regexp.MustCompile(regex)
}

var pathlistCmd = &cobra.Command{
	Use:   "pathlist",
	Short: "Work with pathlist files",
}

func pathlistScanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "scan",
		Short: "Scan a core file for interesting paths",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			verbose, _ := cmd.PersistentFlags().GetBool("verbose")
			reader, openErr := os.Open(args[0])
			r := bufio.NewReader(reader)
			if openErr != nil {
				fmt.Fprintln(os.Stderr, openErr)
				os.Exit(1)
			}
			matches := make(map[string]struct{})
			offset := int64(0)
			for {  // Inspired from https://go.dev/play/p/aPrAW7XGHi

				m := pattern.FindReaderIndex(r)
				if m == nil {
					break
				}

				offset += int64(m[0])  // Add the offset.
				reader.Seek(offset, 0)  // Seek to the match.
				r.Reset(reader)  // Flush buffer so buffered reader is looking where we want it to.

				b := make([]byte, m[1]-m[0])
				io.ReadFull(r, b)
				s := string(b[:])
				if verbose {
					fmt.Fprintf(os.Stderr, "offset=%v s=%v\n", offset, s)
				}

				matches[s] = struct{}{}  // Set membership
				offset += int64(m[1] - m[0])  // Keep record of where we think the reader is looking next.
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
