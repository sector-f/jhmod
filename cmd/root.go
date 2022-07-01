package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "jhmod",
	Short: "Manipulate nvc files",
	Long:  `based off jh_extract.py`,
}

func init() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(extractCmd())
	rootCmd.AddCommand(pathlistCmd)
	rootCmd.AddCommand(createCommand())
	rootCmd.AddCommand(saveCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
