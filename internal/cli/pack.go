package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/evaneliasyoung/phex/internal/pack"
)

func RunPack(inputDir, outputDir, packName string, maxSize, padding int) error {
	files, err := filepath.Glob(filepath.Join(inputDir, "*.png"))
	if err != nil {
		return fmt.Errorf("failed to list input directory: %w", err)
	}

	numSprites := len(files)
	if numSprites == 0 {
		return fmt.Errorf("no PNG files found in %s", inputDir)
	}

	fmt.Printf("[info] found %d textures\n", numSprites)

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	sprites, err := pack.LoadAndTrim(files)
	if err != nil {
		return fmt.Errorf("failed to load input images: %w", err)
	}

	deduped, aliasMap := pack.DedupeWithMap(sprites)

	packed, sheets := pack.PackSprites(deduped, maxSize, padding)

	fmt.Printf("[info] writing to %s\n", filepath.Join(outputDir, packName))
	if err := pack.SaveSheets(packed, sheets, packName, outputDir); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}
	if err := pack.SaveJSON(sprites, packed, aliasMap, sheets, packName, outputDir); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	fmt.Printf("[info] packed %d sprites into %d sheets", len(sprites), len(sheets))

	return nil
}
