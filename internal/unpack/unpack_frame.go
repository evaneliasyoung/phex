package unpack

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	_ "golang.org/x/image/webp"

	"github.com/evaneliasyoung/phex/internal/phaser"
)

func (unpacker Unpacker) UnpackFrame(fr phaser.Frame, img image.Image) error {
	spriteSize := fr.SourceSize.Rect()
	sprite := image.NewRGBA(spriteSize)

	destFrame := fr.SpriteSourceSize.Rect()
	sourceFrame := fr.Frame.Rect()

	draw.Draw(sprite, destFrame, img, sourceFrame.Min, draw.Src)

	encoder := png.Encoder{CompressionLevel: png.DefaultCompression}

	outputPath := filepath.Join(unpacker.OutputDir, fr.FileName+".png")

	if strings.Contains(fr.FileName, "/") {
		parts := strings.Split(fr.FileName, "/")
		subDir := filepath.Join(unpacker.OutputDir, filepath.Join(parts[:len(parts)-1]...))

		if err := os.MkdirAll(subDir, 0o755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to open output file: %w", err)
	}

	if err = encoder.Encode(outputFile, sprite); err != nil {
		return fmt.Errorf("failed to encode sprite as png: %w", err)
	}

	if err = outputFile.Close(); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}
