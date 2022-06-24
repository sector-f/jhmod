package cmd

import 	(
	"fmt"
	"os"
	"regexp"
	"bufio"
	"io"

	"github.com/spf13/cobra"
)

const (
	regex = "data/[A-Za-z0-9_./]+\\.[a-zA-Z0-9]+"
)

var pattern regexp.Regexp

func init() {
	pathlistCmd.AddCommand(pathlistScanCmd)
	pattern = *regexp.MustCompile(regex)
}

var pathlistCmd = &cobra.Command{
	Use:   "pathlist",
	Short: "Work with pathlist files",
}

var pathlistScanCmd = &cobra.Command{
	Use: "scan",
	Short: "Scan a core file for interesting paths",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		reader, openErr := os.Open(args[0])
		r := bufio.NewReader(reader)
		if openErr != nil {
			fmt.Fprintln(os.Stderr, openErr)
			os.Exit(1)
		}
		matches := make(map[string]bool)
		// Borrowed from https://go.dev/play/p/aPrAW7XGHi
		for {
			m := pattern.FindReaderIndex(r)
			if m == nil {
				break
			}
			fmt.Fprintln(os.Stderr, m[0])
			b := make([]byte, m[1]-m[0])
			reader.Seek(int64(m[0]), 0)
			io.ReadFull(reader, b)
			matches[string(b[:])] = true
			r.Reset(reader)  // XXX add a comment explaining what this does.
		}
		fmt.Printf("Found %v matches.", len(matches))
	},
}
