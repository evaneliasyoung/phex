package pack_test

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"

	"github.com/evaneliasyoung/phex/internal/pack"
	"github.com/evaneliasyoung/phex/internal/phaser"
	_ "golang.org/x/image/webp"
)

func TestSaveSheetsPreservesTransparentRGB(t *testing.T) {
	dir := t.TempDir()

	trimmed := image.NewNRGBA(image.Rect(0, 0, 2, 1))
	trimmed.SetNRGBA(0, 0, color.NRGBA{R: 255, A: 0})
	trimmed.SetNRGBA(1, 0, color.NRGBA{G: 255, A: 255})

	sprite := &pack.Sprite{
		Name:       "sprite",
		FullSize:   trimmed.Bounds(),
		Trimmed:    trimmed,
		TrimBounds: trimmed.Bounds(),
	}
	packed := []*pack.PackedSprite{{
		Sprite:     sprite,
		SheetIndex: 0,
		Position:   image.Pt(0, 0),
	}}
	sheets := []*pack.Sheet{{
		Size:    phaser.Size{W: 2, H: 1},
		Sprites: packed,
	}}

	if err := pack.SaveSheets(packed, sheets, "atlas", dir); err != nil {
		t.Fatalf("SaveSheets returned error: %v", err)
	}

	f, err := os.Open(filepath.Join(dir, "atlas-0.webp"))
	if err != nil {
		t.Fatalf("failed to open saved sheet: %v", err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			t.Fatalf("failed to close saved sheet: %v", cerr)
		}
	}()

	img, _, err := image.Decode(f)
	if err != nil {
		t.Fatalf("failed to decode saved sheet: %v", err)
	}

	got := color.NRGBAModel.Convert(img.At(0, 0)).(color.NRGBA)
	want := color.NRGBA{R: 255, A: 0}
	if got != want {
		t.Fatalf("transparent pixel mismatch: got %#v, want %#v", got, want)
	}
}
