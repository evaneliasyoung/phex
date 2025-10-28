package cmd

import (
	"fmt"
	"os"

	"github.com/evaneliasyoung/phex/internal"
	"github.com/spf13/cobra"
)

var outputDir string
var showVersion bool

var rootCmd = &cobra.Command{
	Use:   "phex",
	Short: "Pack and unpack Phaser 3 sprite atlases",
	Long: `Phex - Phaser Texture Manager

Pack individual PNGs into WebP sprite sheets and unpack existing atlases
for Phaser 3. Includes transparency trimming, dupl,icate-sprite
deduplication, and optimal packing.`,
	Run: func(cmd *cobra.Command, args []string) {
		if showVersion {
			_, err := fmt.Print(internal.Version)
			cobra.CheckErr(err)
		} else {
			err := cmd.Help()
			cobra.CheckErr(err)
		}
	},
}

func init() {
	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "version of phex")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
