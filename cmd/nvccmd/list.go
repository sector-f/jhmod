package nvccmd

import (
	"fmt"
	"os"

	"github.com/sector-f/jhmod/nvc"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list FILE",
	Short: "Manipulate nvc files",
	Long:  `based off jh_extract.py`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		reader, openErr := os.Open(args[0])
		if openErr != nil {
			fmt.Fprintln(os.Stderr, openErr)
			os.Exit(1)
		}
		defer reader.Close()

		archive, listErr := nvc.Parse(reader)
		if listErr != nil {
			fmt.Fprintln(os.Stderr, listErr)
			os.Exit(1)
		}

		for _, hash := range archive.EntryOrder {
			entry, _ := archive.Entries[hash]
			fmt.Println(entry)
		}
	},
}
