package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "jh_extract",
	Short: "Manipulate nvc files",
	Long:  `based off jh_extract.py`,
}

func init() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(extractCmd())
	rootCmd.AddCommand(pathlistCmd)
	rootCmd.AddCommand(createCommand())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
