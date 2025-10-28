package cli

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	"github.com/evaneliasyoung/phex/internal/pack"
)

func isSupportedFile(path string, d fs.DirEntry) (bool, error) {
	if d.IsDir() {
		return false, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer file.Close()

	head := make([]byte, 261)
	_, err = file.Read(head)
	if err != nil {
		return false, err
	}

	return pack.IsSupportedImage(head), nil
}

func findImageFiles(root string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		isSupported, err := isSupportedFile(path, d)
		if err != nil {
			return err
		}

		if isSupported {
			files = append(files, path)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Strings(files)
	return files, nil
}

func RunPack(inputDir, outputDir, packName string, maxSize, padding int) error {
	files, err := findImageFiles(inputDir)
	if err != nil {
		return fmt.Errorf("failed to list input directory: %w", err)
	}

	numSprites := len(files)
	if numSprites == 0 {
		return fmt.Errorf("no image files found in %s", inputDir)
	}

	fmt.Printf("[info] found %d textures\n", numSprites)

	if outputDir == "" {
		outputDir = filepath.Join(filepath.Dir(inputDir), packName)
	}

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	sprites, err := pack.LoadAndTrim(inputDir, files)
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
