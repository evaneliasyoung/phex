package pack_test

import (
	"image"
	"testing"

	"github.com/evaneliasyoung/phex/internal/pack"
)

func makeSprite(name string, w, h int) *pack.Sprite {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	return &pack.Sprite{
		Name:     name,
		FullSize: img.Bounds(),
		Trimmed:  img,
	}
}

func TestPackSpritesShortSideHeuristic(t *testing.T) {
	sprites := []*pack.Sprite{
		makeSprite("strip", 4, 1),
		makeSprite("square", 2, 2),
		makeSprite("bar", 2, 1),
	}

	packed, sheets, err := pack.PackSprites(sprites, 4, 0)
	if err != nil {
		t.Fatalf("pack.PackSprites returned error: %v", err)
	}

	if len(sheets) != 1 {
		t.Fatalf("expected a single sheet, got %d", len(sheets))
	}
	if sheets[0].W != 4 || sheets[0].H != 4 {
		t.Fatalf("unexpected sheet size %dx%d", sheets[0].W, sheets[0].H)
	}

	placements := make(map[string]struct {
		idx int
		pos image.Point
	})
	for _, p := range packed {
		placements[p.Sprite.Name] = struct {
			idx int
			pos image.Point
		}{idx: p.SheetIndex, pos: p.Position}
	}

	for name, want := range map[string]struct {
		idx int
		pos image.Point
	}{
		"strip":  {idx: 0, pos: image.Pt(0, 0)},
		"square": {idx: 0, pos: image.Pt(0, 1)},
		"bar":    {idx: 0, pos: image.Pt(0, 3)},
	} {
		got, ok := placements[name]
		if !ok {
			t.Fatalf("missing placement for %q", name)
		}
		if got.idx != want.idx {
			t.Fatalf("%s placed on sheet %d, want %d", name, got.idx, want.idx)
		}
		if got.pos != want.pos {
			t.Fatalf("%s positioned at %v, want %v", name, got.pos, want.pos)
		}
	}
}

func TestPackSpritesAllocatesNewSheet(t *testing.T) {
	sprites := []*pack.Sprite{
		makeSprite("strip", 4, 1),
		makeSprite("large", 3, 3),
		makeSprite("extra", 2, 2),
	}

	packed, sheets, err := pack.PackSprites(sprites, 4, 0)
	if err != nil {
		t.Fatalf("pack.PackSprites returned error: %v", err)
	}

	if len(sheets) != 2 {
		t.Fatalf("expected two sheets, got %d", len(sheets))
	}
	if sheets[0].W != 4 || sheets[0].H != 4 {
		t.Fatalf("unexpected first sheet size %dx%d", sheets[0].W, sheets[0].H)
	}
	if sheets[1].W != 2 || sheets[1].H != 2 {
		t.Fatalf("unexpected second sheet size %dx%d", sheets[1].W, sheets[1].H)
	}

	placements := make(map[string]struct {
		idx int
		pos image.Point
	})
	for _, p := range packed {
		placements[p.Sprite.Name] = struct {
			idx int
			pos image.Point
		}{idx: p.SheetIndex, pos: p.Position}
	}

	expectations := map[string]struct {
		idx int
		pos image.Point
	}{
		"strip": {idx: 0, pos: image.Pt(0, 0)},
		"large": {idx: 0, pos: image.Pt(0, 1)},
		"extra": {idx: 1, pos: image.Pt(0, 0)},
	}

	for name, want := range expectations {
		got, ok := placements[name]
		if !ok {
			t.Fatalf("missing placement for %q", name)
		}
		if got.idx != want.idx {
			t.Fatalf("%s placed on sheet %d, want %d", name, got.idx, want.idx)
		}
		if got.pos != want.pos {
			t.Fatalf("%s positioned at %v, want %v", name, got.pos, want.pos)
		}
	}
}
