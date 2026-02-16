package pack

import (
	"image"
	"testing"
)

func makeSprite(name string, w, h int) *Sprite {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	return &Sprite{
		Name:       name,
		FullSize:   img.Bounds(),
		Trimmed:    img,
		TrimBounds: img.Bounds(),
	}
}

func TestPackSpritesProducesValidPacking(t *testing.T) {
	sprites := []*Sprite{
		makeSprite("strip", 4, 1),
		makeSprite("square", 2, 2),
		makeSprite("bar", 2, 1),
	}

	packed, sheets, err := PackSprites(sprites, 4, 0)
	if err != nil {
		t.Fatalf("PackSprites returned error: %v", err)
	}
	if len(sheets) != 1 {
		t.Fatalf("expected a single sheet, got %d", len(sheets))
	}

	assertValidPacking(t, packed, sheets)
	if gotArea := scorePacking(sheets).totalArea; gotArea > 16 {
		t.Fatalf("expected packed area <= 16, got %d", gotArea)
	}
}

func TestPackSpritesAllocatesMultipleSheetsWhenRequired(t *testing.T) {
	sprites := []*Sprite{
		makeSprite("strip", 4, 1),
		makeSprite("large", 3, 3),
		makeSprite("extra", 2, 2),
	}

	packed, sheets, err := PackSprites(sprites, 4, 0)
	if err != nil {
		t.Fatalf("PackSprites returned error: %v", err)
	}
	if len(sheets) < 2 {
		t.Fatalf("expected at least two sheets, got %d", len(sheets))
	}

	assertValidPacking(t, packed, sheets)
}

func assertValidPacking(t *testing.T, packed []*PackedSprite, sheets []*Sheet) {
	t.Helper()

	if len(packed) == 0 {
		return
	}

	rectsBySheet := make(map[int][]image.Rectangle)
	seen := make(map[string]bool)

	for _, ps := range packed {
		if ps.SheetIndex < 0 || ps.SheetIndex >= len(sheets) {
			t.Fatalf("sprite %q placed on invalid sheet index %d", ps.Sprite.Name, ps.SheetIndex)
		}
		if seen[ps.Sprite.Name] {
			t.Fatalf("sprite %q appears more than once", ps.Sprite.Name)
		}
		seen[ps.Sprite.Name] = true

		w, h := spriteSize(ps.Sprite)
		r := image.Rect(ps.Position.X, ps.Position.Y, ps.Position.X+w, ps.Position.Y+h)
		sh := sheets[ps.SheetIndex]
		if r.Max.X > sh.W || r.Max.Y > sh.H || r.Min.X < 0 || r.Min.Y < 0 {
			t.Fatalf("sprite %q rect %v is out of sheet bounds %dx%d", ps.Sprite.Name, r, sh.W, sh.H)
		}
		rectsBySheet[ps.SheetIndex] = append(rectsBySheet[ps.SheetIndex], r)
	}

	for sheetIdx, rects := range rectsBySheet {
		for i := 0; i < len(rects); i++ {
			for j := i + 1; j < len(rects); j++ {
				if rects[i].Overlaps(rects[j]) {
					t.Fatalf("sheet %d has overlapping rects %v and %v", sheetIdx, rects[i], rects[j])
				}
			}
		}
	}
}
