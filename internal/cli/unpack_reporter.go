package cli

import (
	"github.com/evaneliasyoung/phex/internal/phaser"
	"github.com/evaneliasyoung/phex/internal/unpack"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

type unpackReporter struct {
	p        *mpb.Progress
	textures map[string]*mpb.Bar
	total    *mpb.Bar
}

func (r *unpackReporter) FrameProcessed(tex phaser.Texture, _ phaser.Frame) {
	if r.p != nil {
		r.textures[tex.Image].Increment()
		r.total.Increment()
	}
}

func (r *unpackReporter) TextureProcessed(phaser.Texture) {}

func (r *unpackReporter) AtlasProcessed() {
	if r.p != nil {
		r.p.Wait()
	}
}

func MakeUnpackReporter(unpacker unpack.Unpacker, totalTextures *int, noProgress bool) *unpackReporter {
	if !IsTTY || noProgress {
		return &unpackReporter{p: nil, textures: nil, total: nil}
	}

	p := mpb.New()
	sheets := make(map[string]*mpb.Bar)

	for _, sh := range unpacker.Textures {
		*totalTextures += len(sh.Frames)

		sheets[sh.Image] = p.AddBar(
			int64(len(sh.Frames)),
			mpb.PrependDecorators(
				decor.Name(sh.Image+" ", decor.WCSyncWidth),
				decor.CountersNoUnit("%d / %d"),
			),
			mpb.AppendDecorators(
				decor.Percentage(),
			),
		)
	}

	total := p.AddBar(
		int64(*totalTextures),
		mpb.PrependDecorators(
			decor.Name("Total ", decor.WCSyncWidth),
			decor.CountersNoUnit("%d / %d"),
		),
		mpb.AppendDecorators(
			decor.Percentage(),
		),
	)

	return &unpackReporter{p: p, textures: sheets, total: total}
}
