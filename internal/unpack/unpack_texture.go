package unpack

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"sync"

	"github.com/evaneliasyoung/phex/internal/phaser"
	_ "golang.org/x/image/webp"
)

func (unpacker Unpacker) UnpackTexture(tex phaser.Texture, reporter ProgressReporter) error {
	sheetPath := filepath.Join(unpacker.InputDir, tex.Image)

	sheetFile, err := os.Open(sheetPath)
	if err != nil {
		return fmt.Errorf("failed to open texture sheet: %w", err)
	}

	img, _, err := image.Decode(sheetFile)
	if err != nil {
		return fmt.Errorf("failed to decode texture sheet: %w", err)
	}
	err = sheetFile.Close()
	if err != nil {
		return fmt.Errorf("failed to close texture sheet: %w", err)
	}

	jobs := make(chan phaser.Frame)
	results := make(chan error, len(tex.Frames))

	var wg sync.WaitGroup

	for range unpacker.Workers {
		wg.Go(func() {
			for fr := range jobs {
				if err := unpacker.UnpackFrame(fr, img); err != nil {
					results <- err
					return
				}
				reporter.FrameProcessed(tex, fr)
				results <- nil
			}
		})
	}

	go func() {
		for _, fr := range tex.Frames {
			jobs <- fr
		}
		close(jobs)
	}()

	wg.Wait()
	close(results)

	for err := range results {
		if err != nil {
			return err
		}
	}

	reporter.TextureProcessed(tex)

	return nil
}
