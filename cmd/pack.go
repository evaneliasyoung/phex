package cmd

import (
	"errors"

	"github.com/evaneliasyoung/phex/internal/cli"
	"github.com/spf13/cobra"
)

var packName string
var maxSize int
var padding int

var packCmd = &cobra.Command{
	Use:   "pack <input-dir>",
	Short: "Build a Phaser 3 sprite atlas from images",
	Long: `Pack a directory of images into a single sprite sheet (atlas) and JSON for Phaser 3.

Input should be a folder containing PNGs (optionally in subfolders). The output includes a
sprite sheet image (e.g., WebP/PNG) and a Phaser 3 JSON describing frames. Phex can trim
transparent borders, deduplicate identical frames, and use an optimal bin-packing layout.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("pack requires a source directory")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := cli.RunPack(args[0], outputDir, packName, maxSize, padding)
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(packCmd)

	packCmd.Flags().StringVarP(&packName, "name", "n", "atlas", "The name of the sprite sheets and the atlas file")
	packCmd.Flags().StringVarP(&outputDir, "output", "o", "./output", "Output directory for atlas and images")
	packCmd.Flags().IntVarP(&maxSize, "maxsize", "m", 2048, "Maximum width/height of output sheets")
	packCmd.Flags().IntVarP(&padding, "padding", "p", 0, "Padding pixels between sprites in the sheet")
}
