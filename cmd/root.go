package cmd

import 	(
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "jh_extract",
	Short: "Manipulate nvc files",
	Long: `based off jh_extract.py`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("durr")
	},
}

func init() {
	fmt.Println("in root")
	rootCmd.AddCommand(listCmd)
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}
