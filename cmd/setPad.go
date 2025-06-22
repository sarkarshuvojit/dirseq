/*
Copyright © 2025 Shuvojit Sarkar <s15sarkar@yahoo.com>
*/
package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"

	dirseq "github.com/sarkarshuvojit/dirseq/pkg"
	"github.com/spf13/cobra"
)

// setPadCmd represents the setPad command
var setPadCmd = &cobra.Command{
	Use:   "set-pad",
	Short: "Set default padding width for the current directory",
	Long: `Stores the default zero-padding width for the sequence number in the current directory.

Once set, running 'dirseq' will automatically pad the output using this width
unless overridden with the --pad flag.

Examples:
  $ dirset set-pad 4
  $ dirseq
  0001

  $ dirseq --pad 2
  01
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Expecting one arg; the number to be set.")
		}

		padding, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		currentPath, err := os.Getwd()
		if err != nil {
			slog.Error("Failed to get current working directory", "error", err)
			os.Exit(1)
		}
		absPath, err := filepath.Abs(currentPath)
		if err != nil {
			slog.Error("Failed to get absolute path", "for", currentPath, "error", err)
			os.Exit(1)
		}

		store := &dirseq.JsonStore{}
		_, err = store.SetupDatabase()
		if err != nil {
			slog.Error("Failed to setup db path", "error", err)
			os.Exit(1)
		}

		store.UpdatePadding(absPath, padding)

		fmt.Printf("Set padding for %s to %d\n", absPath, padding)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(setPadCmd)
}
