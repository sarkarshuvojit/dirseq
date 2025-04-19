/*
Copyright © 2025 Shuvojit Sarkar <s15sarkar@yahoo.com>
*/
package cmd

import (
	"os"

	dirseq "github.com/sarkarshuvojit/dirseq/pkg"
	"github.com/spf13/cobra"
)

var padding int


var rootCmd = &cobra.Command{
	Use:   "dirseq",
	Short: "Print the next sequence number for the current directory",
	Long: `Outputs the next sequence number associated with the current working directory.

Each time 'dirseq' is run from a directory, it increments and returns the next number
in a sequence specific to that path. This is useful for numbering files, directories,
or logs consistently.

Use the --pad flag to left-pad the number with zeros to a specific width.

Examples:
  $ dirseq
  42

  $ dirseq --pad 4
  0042
  `,
	Run: func(cmd *cobra.Command, args []string) {
		dirseq.Show(padding)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().IntVar(&padding, "p", 0, "Pad the output with leading zeros to reach the specified length (e.g., -p 4 => 0055)")
}


