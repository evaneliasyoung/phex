package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/evaneliasyoung/phex/internal"
	"github.com/spf13/cobra"
)

var outputDir string
var showVersion int

var rootCmd = &cobra.Command{
	Use:   "phex",
	Short: "Pack and unpack Phaser 3 sprite atlases",
	Long: `Phex - Phaser Texture Manager

Pack individual images into WebP sprite sheets and unpack existing atlases
for Phaser 3. Includes transparency trimming, duplicate-sprite
deduplication, and optimal packing.`,
	Run: func(cmd *cobra.Command, args []string) {
		switch showVersion {
		case 0:
			err := cmd.Help()
			cobra.CheckErr(err)
		case 1:
			_, err := fmt.Print(internal.Version)
			cobra.CheckErr(err)
		case 2:
			_, err := fmt.Printf("%s-%s", internal.Version, runtime.GOOS)
			cobra.CheckErr(err)
		case 3:
			_, err := fmt.Printf("%s-%s-%s", internal.Version, runtime.GOOS, runtime.GOARCH)
			cobra.CheckErr(err)
		}
	},
}

func init() {
	rootCmd.Flags().CountVarP(&showVersion, "version", "v", "version of phex")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
