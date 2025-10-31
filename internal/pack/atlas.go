package pack

import (
	"fmt"
	"image"
	"sort"

	"github.com/evaneliasyoung/phex/internal/phaser"
)

type PackedSprite struct {
	Sprite     *Sprite
	SheetIndex int
	Position   image.Point
}

type Sheet struct {
	phaser.Size
	Sprites []*PackedSprite
}

func contains(a, b image.Rectangle) bool {
	return a.Min.X <= b.Min.X && a.Min.Y <= b.Min.Y && a.Max.X >= b.Max.X && a.Max.Y >= b.Max.Y
}

func PackSprites(sprites []*Sprite, maxSize, padding int) ([]*PackedSprite, []*Sheet, error) {
	sorted := make([]*Sprite, len(sprites))
	copy(sorted, sprites)
	sort.Slice(sorted, func(i, j int) bool {
		iw, ih := sorted[i].Trimmed.Bounds().Dx(), sorted[i].Trimmed.Bounds().Dy()
		jw, jh := sorted[j].Trimmed.Bounds().Dx(), sorted[j].Trimmed.Bounds().Dy()
		imax, jmax := max(iw, ih), max(jw, jh)
		if imax != jmax {
			return imax > jmax
		}
		ia, ja := iw*ih, jw*jh
		return ia > ja
	})

	return packMaxRects(sorted, maxSize, padding)
}

func packMaxRects(sprites []*Sprite, maxSize, padding int) ([]*PackedSprite, []*Sheet, error) {
	var sheets []*Sheet

	current := &Sheet{}
	bin := newMaxRectsBin(maxSize, maxSize)
	for _, s := range sprites {
		spriteWidth := s.Trimmed.Bounds().Dx()
		spriteHeight := s.Trimmed.Bounds().Dy()

		if spriteWidth+padding > maxSize || spriteHeight+padding > maxSize {
			return nil, nil, fmt.Errorf("sprite %q (%dx%d) exceeds max sheet size %d with padding %d", s.Name, spriteWidth, spriteHeight, maxSize, padding)
		}

		place := bin.insert(spriteWidth+padding, spriteHeight+padding)
		if !place.ok {
			finalizeSheetSize(current)
			sheets = append(sheets, current)
			current = &Sheet{}
			bin = newMaxRectsBin(maxSize, maxSize)
			place = bin.insert(spriteWidth+padding, spriteHeight+padding)
			if !place.ok {
				return nil, nil, fmt.Errorf("failed to place sprite %q (%dx%d) into a new sheet of size %dx%d", s.Name, spriteWidth, spriteHeight, maxSize, maxSize)
			}
		}
		rot := place.rot
		w := spriteWidth
		h := spriteHeight
		px, py := place.X, place.Y
		ps := &PackedSprite{Sprite: s, SheetIndex: len(sheets), Position: image.Pt(px, py)}
		current.Sprites = append(current.Sprites, ps)
		if !rot {
			current.W = max(current.W, px+w)
			current.H = max(current.H, py+h)
		} else {
			current.W = max(current.W, px+h)
			current.H = max(current.H, py+w)
		}
	}
	finalizeSheetSize(current)
	sheets = append(sheets, current)
	var out []*PackedSprite
	for _, sh := range sheets {
		out = append(out, sh.Sprites...)
	}
	return out, sheets, nil
}

func finalizeSheetSize(sh *Sheet) {
	sh.W = max(sh.W, 1)
	sh.H = max(sh.H, 1)
}
