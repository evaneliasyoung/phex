package cmd

import (
	"errors"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/evaneliasyoung/phex/internal/cli"
	"github.com/spf13/cobra"
)

var workers int
var noProgress bool

var unpackCmd = &cobra.Command{
	Use:   "unpack <atlas>",
	Short: "Extract frames from a Phaser atlas",
	Long: `Unpack a Phaser 3 sprite atlas into individual image files.

Provide the path to a Phaser-compatible atlas JSON (e.g. TexturePacker or Phex output).
This writes each frame to an output directory, preserving any subfolder structure when available.
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 || strings.ToLower(filepath.Ext(args[0])) != ".json" {
			return errors.New("unpack requires an atlas.json file")
		}
		if workers <= 0 {
			return errors.New("workers must be greater than or equal to 1")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := cli.RunUnpack(args[0], outputDir, workers, noProgress)
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(unpackCmd)

	defaultWorkers := min(2*runtime.NumCPU(), 32)

	unpackCmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory")
	unpackCmd.Flags().IntVarP(&workers, "workers", "w", defaultWorkers, "Number of concurrent workers")
	unpackCmd.Flags().BoolVar(&noProgress, "no-progress", !cli.IsTTY, "Disable progress bars")
}
