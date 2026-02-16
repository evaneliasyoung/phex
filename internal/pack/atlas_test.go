package pack_test

import (
	"image"
	"strings"
	"testing"

	"github.com/evaneliasyoung/phex/internal/pack"
)

func newTestSprite(name string, w, h int) *pack.Sprite {
	trimmed := image.NewNRGBA(image.Rect(0, 0, w, h))
	return &pack.Sprite{
		Name:       name,
		FullSize:   trimmed.Bounds(),
		Trimmed:    trimmed,
		TrimBounds: trimmed.Bounds(),
	}
}

func TestPackSpritesSpriteTooLarge(t *testing.T) {
	sprites := []*pack.Sprite{newTestSprite("large", 10, 10)}
	if _, _, err := pack.PackSprites(sprites, 10, 1); err == nil {
		t.Fatalf("expected error when sprite exceeds max sheet size")
	} else if !strings.Contains(err.Error(), "exceeds max sheet size") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPackSpritesCreatesMultipleSheets(t *testing.T) {
	sprites := []*pack.Sprite{
		newTestSprite("first", 10, 10),
		newTestSprite("second", 10, 10),
	}

	packed, sheets, err := pack.PackSprites(sprites, 10, 0)
	if err != nil {
		t.Fatalf("pack.PackSprites returned error: %v", err)
	}

	if len(sheets) != 2 {
		t.Fatalf("expected 2 sheets, got %d", len(sheets))
	}

	for i, sh := range sheets {
		if sh.W != 10 || sh.H != 10 {
			t.Fatalf("sheet %d has dimensions %dx%d, want 10x10", i, sh.W, sh.H)
		}
	}

	if len(packed) != 2 {
		t.Fatalf("expected 2 packed sprites, got %d", len(packed))
	}
}

func TestPackSpritesEmptySpritesProducesUnitSheet(t *testing.T) {
	packed, sheets, err := pack.PackSprites(nil, 10, 0)
	if err != nil {
		t.Fatalf("pack.PackSprites returned error: %v", err)
	}

	if len(packed) != 0 {
		t.Fatalf("expected no packed sprites, got %d", len(packed))
	}

	if len(sheets) != 1 {
		t.Fatalf("expected 1 sheet, got %d", len(sheets))
	}

	if sheets[0].W != 1 || sheets[0].H != 1 {
		t.Fatalf("expected sheet size 1x1, got %dx%d", sheets[0].W, sheets[0].H)
	}
}
