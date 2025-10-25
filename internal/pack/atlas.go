package pack

import (
	"image"
	"sort"

	"github.com/evaneliasyoung/phex/internal/phaser"
)

type PackedSprite struct {
	Sprite     *Sprite
	SheetIndex int
	Position   image.Point
	Rotated    bool
}

type Sheet struct {
	phaser.Size
	Sprites []*PackedSprite
}

func contains(a, b image.Rectangle) bool {
	return a.Min.X <= b.Min.X && a.Min.Y <= b.Min.Y && a.Max.X >= b.Max.X && a.Max.Y >= b.Max.Y
}

func PackSprites(sprites []*Sprite, maxSize, padding int, allowRotate bool) ([]*PackedSprite, []*Sheet) {
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

	return packMaxRects(sorted, maxSize, padding, allowRotate)
}

func packMaxRects(sprites []*Sprite, maxSize, padding int, allowRotate bool) ([]*PackedSprite, []*Sheet) {
	var sheets []*Sheet

	current := &Sheet{}
	bin := newMaxRectsBin(maxSize, maxSize, allowRotate)
	for _, s := range sprites {
		place := bin.insert(s.Trimmed.Bounds().Dx()+padding, s.Trimmed.Bounds().Dy()+padding)
		if !place.ok {
			finalizeSheetSize(current)
			sheets = append(sheets, current)
			current = &Sheet{}
			bin = newMaxRectsBin(maxSize, maxSize, allowRotate)
			place = bin.insert(s.Trimmed.Bounds().Dx()+padding, s.Trimmed.Bounds().Dy()+padding)
			if !place.ok {
				place = bin.insert(s.Trimmed.Bounds().Dx(), s.Trimmed.Bounds().Dy())
				if !place.ok {
					bin = newMaxRectsBin(max(s.Trimmed.Bounds().Dx(), maxSize), max(s.Trimmed.Bounds().Dy(), maxSize), allowRotate)
					place = bin.insert(s.Trimmed.Bounds().Dx(), s.Trimmed.Bounds().Dy())
				}
			}
		}
		rot := place.rot
		w := s.Trimmed.Bounds().Dx()
		h := s.Trimmed.Bounds().Dy()
		px, py := place.X, place.Y
		ps := &PackedSprite{Sprite: s, SheetIndex: len(sheets), Position: image.Pt(px, py), Rotated: rot}
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
	return out, sheets
}

func finalizeSheetSize(sh *Sheet) {
	sh.W = max(sh.W, 1)
	sh.H = max(sh.H, 1)
}
