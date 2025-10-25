package pack

import (
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"os"
	"path/filepath"

	"github.com/evaneliasyoung/phex/internal"
	"github.com/evaneliasyoung/phex/internal/phaser"
	"github.com/gen2brain/webp"
)

func SaveSheets(packed []*PackedSprite, sheets []*Sheet, packName, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	for i, sh := range sheets {
		w := max(sh.W, 1)
		h := max(sh.H, 1)
		canvas := image.NewNRGBA(image.Rect(0, 0, w, h))
		for _, ps := range sh.Sprites {
			src := ps.Sprite.Trimmed
			dst := image.Rect(ps.Position.X, ps.Position.Y, ps.Position.X+src.Bounds().Dx(), ps.Position.Y+src.Bounds().Dy())
			draw.Draw(canvas, dst, src, image.Point{}, draw.Src)
		}
		outPath := filepath.Join(outputDir, fmt.Sprintf("%s-%d.webp", packName, i))
		f, err := os.Create(outPath)
		if err != nil {
			return fmt.Errorf("failed to open output file: %w", err)
		}

		if err := webp.Encode(f, canvas, webp.Options{Lossless: true}); err != nil {
			_ = f.Close()
			return fmt.Errorf("failed to write output file: %w", err)
		}
		_ = f.Close()
	}
	return nil
}

func SaveJSON(original []*Sprite, packed []*PackedSprite, aliasMap map[string]string, sheets []*Sheet, packName, outputDir string) error {
	frameMap := make(map[string]*PackedSprite)
	for _, ps := range packed {
		frameMap[ps.Sprite.Name] = ps
	}

	var atlas phaser.Atlas
	atlas.Meta = map[string]string{
		"app":     "https://github.com/evaneliasyoung/phex",
		"version": internal.Version,
	}
	for i, sh := range sheets {
		tex := phaser.Texture{
			Image:  fmt.Sprintf("%s-%d.webp", packName, i),
			Format: "RGBA8888",
			Size:   sh.Size,
			Scale:  1,
		}
		for _, orig := range original {
			canName, ok := aliasMap[orig.Name]
			if !ok {
				continue
			}
			ps := frameMap[canName]
			if ps == nil || ps.SheetIndex != i {
				continue
			}
			trim := orig.TrimBounds
			fw, fh := trim.Dx(), trim.Dy()
			tex.Frames = append(tex.Frames, phaser.Frame{
				FileName:         orig.Name,
				Rotated:          false,
				Trimmed:          orig.WasTrimmed,
				SourceSize:       phaser.Size{W: orig.FullSize.Dx(), H: orig.FullSize.Dy()},
				SpriteSourceSize: phaser.Rect{X: trim.Min.X, Y: trim.Min.Y, W: trim.Dx(), H: trim.Dy()},
				Frame:            phaser.Rect{X: ps.Position.X, Y: ps.Position.Y, W: fw, H: fh},
			})
		}
		atlas.Textures = append(atlas.Textures, tex)
	}
	data, err := json.MarshalIndent(atlas, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to generate JSON file: %w", err)
	}
	return os.WriteFile(filepath.Join(outputDir, fmt.Sprintf("%s.json", packName)), data, 0o644)
}
