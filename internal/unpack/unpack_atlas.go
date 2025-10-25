package unpack

import (
	"fmt"
	"os"
	"sync"

	"github.com/evaneliasyoung/phex/internal/phaser"
)

func (unpacker Unpacker) UnpackAtlas(reporter ProgressReporter) error {
	if err := os.MkdirAll(unpacker.OutputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var firstErr error

	for _, sh := range unpacker.Textures {
		wg.Add(1)

		go func(tex phaser.Texture, pr ProgressReporter) {
			defer wg.Done()
			if err := unpacker.UnpackTexture(tex, pr); err != nil {
				mu.Lock()
				if firstErr == nil {
					firstErr = err
				}
				mu.Unlock()
			}
		}(sh, reporter)
	}

	wg.Wait()

	reporter.AtlasProcessed()

	if firstErr != nil {
		return firstErr
	}

	return nil
}
