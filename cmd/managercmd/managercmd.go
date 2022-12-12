package managercmd

import (
	"github.com/sector-f/jhmod/gui"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "manager",
		Run: func(cmd *cobra.Command, args []string) {
			gui.Run()
		},
	}

	// Flag go here

	return cmd
}
