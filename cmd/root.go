package cmd

import (
	"fmt"
	"os"

	"github.com/sector-f/jhmod/cmd/nvccmd"
	"github.com/sector-f/jhmod/cmd/savecmd"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "jhmod",
	Short: "Manipulate nvc files",
	Long:  `based off jh_extract.py`,
}

func init() {
	rootCmd.AddCommand(nvccmd.Cmd())
	rootCmd.AddCommand(savecmd.Cmd())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
