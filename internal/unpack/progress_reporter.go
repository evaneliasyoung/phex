package unpack

import "github.com/evaneliasyoung/phex/internal/phaser"

type ProgressReporter interface {
	FrameProcessed(texture phaser.Texture, frame phaser.Frame)
	TextureProcessed(texture phaser.Texture)
	AtlasProcessed()
}

type NoProgressReporter struct{}

func (*NoProgressReporter) FrameProcessed(phaser.Texture, phaser.Frame) {}
func (*NoProgressReporter) TextureProcessed(phaser.Texture)             {}
func (*NoProgressReporter) AtlasProcessed()                             {}
