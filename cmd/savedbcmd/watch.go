package savedbcmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/sector-f/jhmod/savefile"
	"github.com/spf13/cobra"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch game directory for save files and automatically back them up",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO:
		// * Allow configuration
		// * Maybe use XDG spec when determining default directory
		//   * First figure out if Steam actually respects XDG in the first place though
		// * Default support for Windows
		//   * No idea how that's supposed to work!
		homeDir, ok := os.LookupEnv("HOME")
		if !ok {
			return errors.New("HOME env var is not set")
		}

		jhDirPath := filepath.Join(homeDir, ".local/share/Steam/steamapps/common/Jupiter Hell/")
		dirEntries, err := os.ReadDir(jhDirPath)
		if err != nil {
			return err
		}

		existingFiles := []string{}
		for _, entry := range dirEntries {
			fullPath := filepath.Join(jhDirPath, entry.Name())

			// Skip everything that isn't a regular file
			if fType := entry.Type(); !fType.IsRegular() {
				continue
			}

			// Skip files with names that don't look like save files
			if _, err := parseSaveName(entry.Name()); err != nil {
				continue
			}

			existingFiles = append(existingFiles, fullPath)
		}

		for _, filename := range existingFiles {
			file, err := os.Open(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error opening %s: %v\n", filename, err)
				continue
			}
			defer file.Close()

			save, err := savefile.Parse(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing %s: %v\n", filename, err)
			}

			fmt.Printf("%s,%s,%s,%d\n", save.PlayerName, save.GameMode, save.CurrentLevel, save.Seed)
		}

		return nil
	},
}

func parseSaveName(filename string) (int, error) {
	if !strings.HasPrefix(filename, "save") {
		return 0, errors.New("filename does not begin with \"save\"")
	}

	n := strings.TrimPrefix(filename, "save")
	return strconv.Atoi(n)
}
