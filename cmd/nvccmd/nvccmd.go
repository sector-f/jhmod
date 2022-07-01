package nvccmd

import (
	"github.com/spf13/cobra"
)

var nvcCmd = &cobra.Command{
	Use:   "nvc",
	Short: "Manipulate nvc files",
}

func init() {
	nvcCmd.AddCommand(listCmd)
	nvcCmd.AddCommand(extractCmd())
	nvcCmd.AddCommand(pathlistCmd)
	nvcCmd.AddCommand(createCommand())
}

func Cmd() *cobra.Command {
	return nvcCmd
}
