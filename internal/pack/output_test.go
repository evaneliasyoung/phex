package pack_test

import (
	"encoding/json"
	"image"
	"image/color"
	"image/draw"
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

func TestSaveOutputFramesLineUpWithMultipleSpritesOnSheet(t *testing.T) {
	sprites := []*pack.Sprite{
		newPatternSprite("red", image.Rect(0, 0, 6, 5), image.Rect(1, 2, 4, 4), color.NRGBA{R: 255, A: 255}),
		newPatternSprite("green", image.Rect(0, 0, 5, 6), image.Rect(2, 1, 4, 5), color.NRGBA{G: 255, A: 255}),
		newPatternSprite("blue", image.Rect(0, 0, 4, 4), image.Rect(0, 0, 2, 3), color.NRGBA{B: 255, A: 255}),
	}

	_, sheets, err := assertSavedAtlasFramesLineUp(t, sprites, 16, 1)
	if err != nil {
		t.Fatalf("PackSprites returned error: %v", err)
	}
	if len(sheets) != 1 {
		t.Fatalf("expected all sprites on one sheet, got %d sheets", len(sheets))
	}
	if len(sheets[0].Sprites) != len(sprites) {
		t.Fatalf("expected %d sprites on sheet, got %d", len(sprites), len(sheets[0].Sprites))
	}
}

func TestSaveOutputFramesLineUpAcrossMultipleSheets(t *testing.T) {
	sprites := []*pack.Sprite{
		newPatternSprite("wide-red", image.Rect(0, 0, 4, 2), image.Rect(0, 0, 4, 2), color.NRGBA{R: 255, A: 255}),
		newPatternSprite("small-green", image.Rect(0, 0, 3, 3), image.Rect(1, 1, 3, 3), color.NRGBA{G: 255, A: 255}),
		newPatternSprite("wide-blue", image.Rect(0, 0, 4, 2), image.Rect(0, 0, 4, 2), color.NRGBA{B: 255, A: 255}),
		newPatternSprite("small-yellow", image.Rect(0, 0, 3, 3), image.Rect(0, 0, 2, 2), color.NRGBA{R: 255, G: 255, A: 255}),
	}

	packed, sheets, err := assertSavedAtlasFramesLineUp(t, sprites, 4, 0)
	if err != nil {
		t.Fatalf("PackSprites returned error: %v", err)
	}
	if len(sheets) < 2 {
		t.Fatalf("expected multiple sheets, got %d", len(sheets))
	}
	if len(packed) != len(sprites) {
		t.Fatalf("expected %d packed sprites, got %d", len(sprites), len(packed))
	}
}

func assertSavedAtlasFramesLineUp(t *testing.T, sprites []*pack.Sprite, maxSize, padding int) ([]*pack.PackedSprite, []*pack.Sheet, error) {
	t.Helper()

	dir := t.TempDir()

	packed, sheets, err := pack.PackSprites(sprites, maxSize, padding)
	if err != nil {
		return nil, nil, err
	}

	aliasMap := map[string]string{}
	for _, sprite := range sprites {
		aliasMap[sprite.Name] = sprite.Name
	}

	if err := pack.SaveSheets(packed, sheets, "atlas", dir); err != nil {
		t.Fatalf("SaveSheets returned error: %v", err)
	}
	if err := pack.SaveJSON(sprites, packed, aliasMap, sheets, "atlas", dir); err != nil {
		t.Fatalf("SaveJSON returned error: %v", err)
	}

	var atlas phaser.Atlas
	data, err := os.ReadFile(filepath.Join(dir, "atlas.json"))
	if err != nil {
		t.Fatalf("failed to read atlas JSON: %v", err)
	}
	if err := json.Unmarshal(data, &atlas); err != nil {
		t.Fatalf("failed to decode atlas JSON: %v", err)
	}
	if len(atlas.Textures) != len(sheets) {
		t.Fatalf("expected %d textures in atlas, got %d", len(sheets), len(atlas.Textures))
	}

	spriteByName := map[string]*pack.Sprite{}
	for _, sprite := range sprites {
		spriteByName[sprite.Name] = sprite
	}

	for _, tex := range atlas.Textures {
		sheetImage := decodeAtlasSheet(t, filepath.Join(dir, tex.Image))

		for _, fr := range tex.Frames {
			sprite := spriteByName[fr.FileName]
			if sprite == nil {
				t.Fatalf("unexpected frame %q", fr.FileName)
			}

			got := image.NewNRGBA(fr.SourceSize.Rect())
			draw.Draw(got, fr.SpriteSourceSize.Rect(), sheetImage, fr.Frame.Rect().Min, draw.Src)

			want := image.NewNRGBA(sprite.FullSize)
			draw.Draw(want, sprite.TrimBounds, sprite.Trimmed, image.Point{}, draw.Src)

			assertNRGBAEqual(t, fr.FileName, got, want)
		}
	}

	return packed, sheets, nil
}

func decodeAtlasSheet(t *testing.T, path string) image.Image {
	t.Helper()

	sheetFile, err := os.Open(path)
	if err != nil {
		t.Fatalf("failed to open atlas sheet: %v", err)
	}
	defer func() {
		if cerr := sheetFile.Close(); cerr != nil {
			t.Fatalf("failed to close atlas sheet: %v", cerr)
		}
	}()

	sheetImage, _, err := image.Decode(sheetFile)
	if err != nil {
		t.Fatalf("failed to decode atlas sheet: %v", err)
	}
	return sheetImage
}

func newPatternSprite(name string, full, trim image.Rectangle, base color.NRGBA) *pack.Sprite {
	trimmed := image.NewNRGBA(image.Rect(0, 0, trim.Dx(), trim.Dy()))
	for y := 0; y < trim.Dy(); y++ {
		for x := 0; x < trim.Dx(); x++ {
			c := base
			c.R += uint8(x)
			c.G += uint8(y)
			c.B += uint8(x + y)
			trimmed.SetNRGBA(x, y, c)
		}
	}

	return &pack.Sprite{
		Name:       name,
		FullSize:   full,
		Trimmed:    trimmed,
		TrimBounds: trim,
		WasTrimmed: !trim.Eq(full),
	}
}

func assertNRGBAEqual(t *testing.T, name string, got, want *image.NRGBA) {
	t.Helper()

	if got.Bounds() != want.Bounds() {
		t.Fatalf("%s bounds mismatch: got %v, want %v", name, got.Bounds(), want.Bounds())
	}

	b := got.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			if got, want := got.NRGBAAt(x, y), want.NRGBAAt(x, y); got != want {
				t.Fatalf("%s pixel mismatch at %d,%d: got %#v, want %#v", name, x, y, got, want)
			}
		}
	}
}
