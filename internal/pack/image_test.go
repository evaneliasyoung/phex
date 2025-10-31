package pack_test

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/evaneliasyoung/phex/internal/pack"
	"github.com/gen2brain/webp"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
)

func TestLoadAndTrimReturnsTrimmedSprite(t *testing.T) {
	dir := t.TempDir()
	img := image.NewNRGBA(image.Rect(0, 0, 4, 4))
	img.SetNRGBA(1, 1, color.NRGBA{R: 255, A: 255})
	img.SetNRGBA(2, 2, color.NRGBA{G: 128, A: 255})

	path := filepath.Join(dir, "sprite.png")
	writePNG(t, path, img)

	sprites, err := pack.LoadAndTrim(dir, []string{path})
	if err != nil {
		t.Fatalf("LoadAndTrim returned error: %v", err)
	}
	if len(sprites) != 1 {
		t.Fatalf("expected 1 sprite, got %d", len(sprites))
	}

	sprite := sprites[0]
	if sprite.Name != "sprite" {
		t.Fatalf("expected name 'sprite', got %q", sprite.Name)
	}

	expectedFull := image.Rect(0, 0, 4, 4)
	if sprite.FullSize != expectedFull {
		t.Fatalf("expected full size %v, got %v", expectedFull, sprite.FullSize)
	}

	expectedTrimBounds := image.Rect(1, 1, 3, 3)
	if sprite.TrimBounds != expectedTrimBounds {
		t.Fatalf("expected trim bounds %v, got %v", expectedTrimBounds, sprite.TrimBounds)
	}

	expectedTrimmedBounds := image.Rect(0, 0, 2, 2)
	if sprite.Trimmed.Bounds() != expectedTrimmedBounds {
		t.Fatalf("expected trimmed bounds %v, got %v", expectedTrimmedBounds, sprite.Trimmed.Bounds())
	}

	if !sprite.WasTrimmed {
		t.Fatalf("expected sprite to be trimmed")
	}

	if got, want := sprite.Hash, hashNRGBA(sprite.Trimmed); got != want {
		t.Fatalf("unexpected hash: got %s, want %s", got, want)
	}
}

func TestLoadAndTrimTransparentImage(t *testing.T) {
	dir := t.TempDir()
	img := image.NewNRGBA(image.Rect(0, 0, 3, 2))

	path := filepath.Join(dir, "empty.png")
	writePNG(t, path, img)

	sprites, err := pack.LoadAndTrim(dir, []string{path})
	if err != nil {
		t.Fatalf("LoadAndTrim returned error: %v", err)
	}
	if len(sprites) != 1 {
		t.Fatalf("expected 1 sprite, got %d", len(sprites))
	}

	sprite := sprites[0]
	if sprite.Name != "empty" {
		t.Fatalf("expected name 'empty', got %q", sprite.Name)
	}

	expectedFull := image.Rect(0, 0, 3, 2)
	if sprite.FullSize != expectedFull {
		t.Fatalf("expected full size %v, got %v", expectedFull, sprite.FullSize)
	}

	expectedTrimBounds := image.Rect(0, 0, 1, 1)
	if sprite.TrimBounds != expectedTrimBounds {
		t.Fatalf("expected trim bounds %v, got %v", expectedTrimBounds, sprite.TrimBounds)
	}

	if sprite.Trimmed.Bounds() != expectedTrimBounds {
		t.Fatalf("expected trimmed bounds %v, got %v", expectedTrimBounds, sprite.Trimmed.Bounds())
	}

	if sprite.WasTrimmed {
		t.Fatalf("expected sprite not to be trimmed")
	}

	if got, want := sprite.Hash, hashNRGBA(sprite.Trimmed); got != want {
		t.Fatalf("unexpected hash: got %s, want %s", got, want)
	}
}

func TestDedupeWithMap(t *testing.T) {
	sprites := []*pack.Sprite{
		{Name: "first", Hash: "hash-1"},
		{Name: "second", Hash: "hash-1"},
		{Name: "third", Hash: "hash-2"},
	}

	unique, alias := pack.DedupeWithMap(sprites)

	if len(unique) != 2 {
		t.Fatalf("expected 2 unique sprites, got %d", len(unique))
	}

	if unique[0].Name != "first" || unique[1].Name != "third" {
		t.Fatalf("unexpected unique order: %#v", unique)
	}

	if alias["first"] != "first" {
		t.Fatalf("expected alias for 'first' to be 'first', got %q", alias["first"])
	}

	if alias["second"] != "first" {
		t.Fatalf("expected alias for 'second' to be 'first', got %q", alias["second"])
	}

	if alias["third"] != "third" {
		t.Fatalf("expected alias for 'third' to be 'third', got %q", alias["third"])
	}
}

func TestIsSupportedImage(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 2, 2))
	img.SetNRGBA(0, 0, color.NRGBA{R: 255, A: 255})

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("failed to encode PNG: %v", err)
	}

	if !pack.IsSupportedImage(buf.Bytes()) {
		t.Fatalf("expected PNG buffer to be supported")
	}

	if err := jpeg.Encode(&buf, img, &jpeg.Options{}); err != nil {
		t.Fatalf("failed to encode JPEG: %v", err)
	}

	if !pack.IsSupportedImage(buf.Bytes()) {
		t.Fatalf("expected JPEG buffer to be supported")
	}

	if err := gif.Encode(&buf, img, &gif.Options{}); err != nil {
		t.Fatalf("failed to encode GIF: %v", err)
	}

	if !pack.IsSupportedImage(buf.Bytes()) {
		t.Fatalf("expected GIF buffer to be supported")
	}

	if err := bmp.Encode(&buf, img); err != nil {
		t.Fatalf("failed to encode BMP: %v", err)
	}

	if !pack.IsSupportedImage(buf.Bytes()) {
		t.Fatalf("expected BMP buffer to be supported")
	}

	if err := tiff.Encode(&buf, img, &tiff.Options{}); err != nil {
		t.Fatalf("failed to encode TIFF: %v", err)
	}

	if !pack.IsSupportedImage(buf.Bytes()) {
		t.Fatalf("expected TIFF buffer to be supported")
	}

	if err := webp.Encode(&buf, img); err != nil {
		t.Fatalf("failed to encode WebP: %v", err)
	}

	if !pack.IsSupportedImage(buf.Bytes()) {
		t.Fatalf("expected WebP buffer to be supported")
	}

	if pack.IsSupportedImage([]byte("not-an-image")) {
		t.Fatalf("expected invalid buffer to be unsupported")
	}
}

func writePNG(t *testing.T, path string, img image.Image) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("failed to create %s: %v", path, err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			t.Fatalf("failed to close %s: %v", path, cerr)
		}
	}()
	if err := png.Encode(f, img); err != nil {
		t.Fatalf("failed to encode %s: %v", path, err)
	}
}

func hashNRGBA(img *image.NRGBA) string {
	h := md5.New()
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		row := img.Pix[(y-b.Min.Y)*img.Stride : (y-b.Min.Y+1)*img.Stride]
		h.Write(row[:b.Dx()*4])
	}
	return hex.EncodeToString(h.Sum(nil))
}
