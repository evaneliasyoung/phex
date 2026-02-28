package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/evaneliasyoung/phex/internal/phaser"
	"github.com/evaneliasyoung/phex/internal/unpack"
)

func RunUnpack(atlasPath, outputDir string, workers int, noProgress bool) error {
	data, err := os.ReadFile(atlasPath)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	var atlas phaser.Atlas
	if err := json.Unmarshal(data, &atlas); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	inputDir := filepath.Dir(atlasPath)
	atlasName := filepath.Base(atlasPath)
	atlasExt := filepath.Ext(atlasName)

	packName := atlasName
	if strings.EqualFold(atlasExt, ".json") {
		packName = strings.TrimSuffix(atlasName, atlasExt)
	}

	if outputDir == "" {
		outputDir = filepath.Join(inputDir, packName)
	}

	unpacker := unpack.Unpacker{
		Atlas:     atlas,
		PackName:  packName,
		InputDir:  inputDir,
		OutputDir: outputDir,
		Workers:   workers,
	}

	numSheets := len(unpacker.Textures)

	fmt.Printf("[info] found %d texture sheets\n", numSheets)
	fmt.Printf("[info] writing to %s\n", unpacker.OutputDir)

	totalTextures := 0
	reporter := MakeUnpackReporter(unpacker, &totalTextures, noProgress)

	if err := unpacker.UnpackAtlas(reporter); err != nil {
		return err
	}

	fmt.Printf("[info] extracted %d sprites from %d texture sheets\n", totalTextures, len(unpacker.Textures))

	return nil
}
