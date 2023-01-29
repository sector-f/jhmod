package savecmd

import (
	"fmt"
	"os"

	"github.com/sector-f/jhmod/savefile"
	"github.com/spf13/cobra"
)

func init() {
	saveCmd.AddCommand(info())

}

var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "Work with save files",
}

func info() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info FILE ...",
		Short: "Show information about a save file",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			debug, _ := cmd.PersistentFlags().GetBool("debug")
			errors := int(0)
			for _, p := range args {
				if debug {
					fmt.Fprintf(os.Stderr, "Reading file %s\n", p)
				}
				save, parseErr := savefile.ParseFile(p)
				if parseErr != nil {
					fmt.Fprintf(os.Stderr, "Failed to parse '%s: %v\n", p, parseErr)
					errors++
					continue
				}
				fmt.Printf("Savefile %s\n", p)
				fmt.Printf("  Game mode: %s\n", save.GameMode)
				fmt.Printf("  Name:      %s\n", save.PlayerName)
				fmt.Printf("  Level:     %s\n", save.CurrentLevel)
				fmt.Printf("  Seed:      %v\n", save.Seed)
				fmt.Println()
			}
			os.Exit(errors)
		},
	}
	cmd.PersistentFlags().BoolP("debug", "d", false, "Show more internal state of the tool.")

	return cmd
}

func Cmd() *cobra.Command {
	return saveCmd
}
