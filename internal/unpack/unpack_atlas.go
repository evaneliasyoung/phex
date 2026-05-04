package unpack

import (
	"fmt"
	"os"
)

func (unpacker Unpacker) UnpackAtlas(reporter ProgressReporter) error {
	if err := os.MkdirAll(unpacker.OutputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	for _, tex := range unpacker.Textures {
		if err := unpacker.UnpackTexture(tex, reporter); err != nil {
			reporter.AtlasProcessed()
			return err
		}
	}

	reporter.AtlasProcessed()

	return nil
}
