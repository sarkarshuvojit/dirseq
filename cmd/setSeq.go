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

// setSeqCmd represents the setSeq command
var setSeqCmd = &cobra.Command{
	Use:   "set-seq",
	Short: "Manually set the sequence number for the current directory",
	Long: `Set a specific starting sequence number for the current directory.

This command maps the given sequence number to the present working directory.
The next time 'dirseq' is called from this directory, it will return the next number
in the sequence (i.e., one higher than the value set here).

Example:
  $ dirset set-seq 55
  $ dirseq
  56`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Expecting one arg; the number to be set.")
		}

		overrideSeq, err := strconv.Atoi(args[0])
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

		store.UpdateSequence(absPath, overrideSeq)

		dirseq.PPrinter.Info(fmt.Sprintf("Set sequence to %d", overrideSeq))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(setSeqCmd)
}
