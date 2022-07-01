package savedbcmd

import "github.com/spf13/cobra"

var savedbCmd = &cobra.Command{
	Use:   "savedb",
	Short: "Interact with the save file database",
}

func Cmd() *cobra.Command {
	savedbCmd.AddCommand(watchCmd)

	return savedbCmd
}
