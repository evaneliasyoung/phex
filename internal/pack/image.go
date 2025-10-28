package pack

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

type Sprite struct {
	Name       string
	FullSize   image.Rectangle
	Trimmed    *image.NRGBA
	TrimBounds image.Rectangle
	WasTrimmed bool
	Hash       string
}

func LoadAndTrim(root string, paths []string) ([]*Sprite, error) {
	var sprites []*Sprite
	for _, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("failed to open %s: %w", path, err)
		}
		img, err := png.Decode(f)
		_ = f.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to decode %s: %w", path, err)
		}

		nrgba := image.NewNRGBA(img.Bounds())
		draw.Draw(nrgba, nrgba.Bounds(), img, img.Bounds().Min, draw.Src)

		trimmed, bounds, wasTrimmed := trimNRGBA(nrgba)
		hash := hashImage(trimmed)

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return nil, fmt.Errorf("failed to determine frame name for %s: %w", path, err)
		}
		rel = filepath.ToSlash(rel)
		ext := filepath.Ext(rel)
		name := strings.TrimSuffix(rel, ext)

		s := &Sprite{
			Name:       name,
			FullSize:   nrgba.Bounds(),
			Trimmed:    trimmed,
			TrimBounds: bounds,
			WasTrimmed: wasTrimmed,
			Hash:       hash,
		}
		sprites = append(sprites, s)
	}
	return sprites, nil
}

func trimNRGBA(img *image.NRGBA) (*image.NRGBA, image.Rectangle, bool) {
	b := img.Bounds()
	minX, minY := b.Max.X, b.Max.Y
	maxX, maxY := b.Min.X-1, b.Min.Y-1

	for y := b.Min.Y; y < b.Max.Y; y++ {
		row := img.Pix[(y-b.Min.Y)*img.Stride : (y-b.Min.Y+1)*img.Stride]
		for x := b.Min.X; x < b.Max.X; x++ {
			a := row[(x-b.Min.X)*4+3]
			if a != 0 {
				minX = min(x, minX)
				maxX = max(x, maxX)
				minY = min(y, minY)
				maxY = max(y, maxY)
			}
		}
	}

	if maxX < minX || maxY < minY {
		out := image.NewNRGBA(image.Rect(0, 0, 1, 1))
		return out, image.Rect(0, 0, 1, 1), false
	}

	trimRect := image.Rect(minX, minY, maxX+1, maxY+1)
	w, h := trimRect.Dx(), trimRect.Dy()
	out := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		copy(out.Pix[y*out.Stride:y*out.Stride+w*4], img.Pix[(y+trimRect.Min.Y-b.Min.Y)*img.Stride+(trimRect.Min.X-b.Min.X)*4:][:w*4])
	}
	return out, trimRect, !trimRect.Eq(img.Bounds())
}

func hashImage(img *image.NRGBA) string {
	h := md5.New()
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		row := img.Pix[(y-b.Min.Y)*img.Stride : (y-b.Min.Y+1)*img.Stride]
		h.Write(row[:(b.Dx() * 4)])
	}
	return hex.EncodeToString(h.Sum(nil))
}

func DedupeWithMap(sprites []*Sprite) ([]*Sprite, map[string]string) {
	seen := make(map[string]*Sprite)
	alias := make(map[string]string)
	var unique []*Sprite

	for _, s := range sprites {
		if can, ok := seen[s.Hash]; ok {
			alias[s.Name] = can.Name
			continue
		}
		seen[s.Hash] = s
		alias[s.Name] = s.Name
		unique = append(unique, s)
	}
	return unique, alias
}
