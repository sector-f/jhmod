package scumcmd

import (
	"fmt"

	"github.com/sector-f/jhmod/gui"
	"github.com/spf13/cobra"
)

func init() {
	scumCmd.AddCommand(list())
	scumCmd.AddCommand(restore())
}

func Cmd() *cobra.Command {
	return scumCmd
}

var scumCmd = &cobra.Command{
	Use:   "scum",
	Short: "Manage save scums",
}

func restore() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restore ID",
		Short: "Restores savescum db entry ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			db := gui.Connect()
			var saveFile gui.StoredSaveFile
			db.First(&saveFile, args[0])
			if err := saveFile.Restore(); err != nil {
				panic(err)
			}
			fmt.Printf("Restored %s to %s\n", saveFile.AbsPath(), saveFile.OriginalBase)
		},
	}
	return cmd
}

func list() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List savescum database entries",
		Run: func(cmd *cobra.Command, args []string) {
			var saveFiles []gui.StoredSaveFile
			gui.Connect().Find(&saveFiles)
			for _, saveFile := range saveFiles {
				fmt.Println(saveFile)
			}
		},
	}
	return cmd
}
