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
	variants := [][]*Sprite{
		sortSprites(sprites, lessByMaxSideThenArea),
		sortSprites(sprites, lessByAreaThenMaxSide),
	}
	heuristics := []maxRectsHeuristic{
		maxRectsShortSideFit,
		maxRectsLongSideFit,
		maxRectsAreaFit,
		maxRectsBottomLeft,
		maxRectsContactPoint,
	}
	strategies := []func([]*Sprite, int, int, maxRectsHeuristic) ([]*PackedSprite, []*Sheet, error){
		packMaxRectsNextFit,
		packMaxRectsBestBin,
		packMaxRectsFillSheet,
	}

	var (
		bestPacked []*PackedSprite
		bestSheets []*Sheet
		bestScore  packingScore
		hasBest    bool
	)

	for _, variant := range variants {
		for _, heuristic := range heuristics {
			for _, strategy := range strategies {
				packed, sheets, err := strategy(variant, maxSize, padding, heuristic)
				if err != nil {
					return nil, nil, err
				}
				score := scorePacking(sheets)
				if !hasBest || betterPackingScore(score, bestScore) {
					bestPacked, bestSheets, bestScore = packed, sheets, score
					hasBest = true
				}
			}
		}
	}

	if !hasBest {
		return nil, nil, fmt.Errorf("failed to pack sprites")
	}
	return bestPacked, bestSheets, nil
}

func packMaxRectsNextFit(sprites []*Sprite, maxSize, padding int, heuristic maxRectsHeuristic) ([]*PackedSprite, []*Sheet, error) {
	var sheets []*Sheet

	current := &Sheet{}
	bin := newMaxRectsBin(maxSize, maxSize)
	for _, s := range sprites {
		spriteWidth, spriteHeight, err := validateSpriteSize(s, maxSize, padding)
		if err != nil {
			return nil, nil, err
		}

		place := bin.insertWithHeuristic(spriteWidth+padding, spriteHeight+padding, heuristic)
		if !place.ok {
			if len(current.Sprites) > 0 {
				finalizeSheetSize(current)
				sheets = append(sheets, current)
			}
			current = &Sheet{}
			bin = newMaxRectsBin(maxSize, maxSize)
			place = bin.insertWithHeuristic(spriteWidth+padding, spriteHeight+padding, heuristic)
			if !place.ok {
				return nil, nil, fmt.Errorf("failed to place sprite %q (%dx%d) into a new sheet of size %dx%d", s.Name, spriteWidth, spriteHeight, maxSize, maxSize)
			}
		}
		w := spriteWidth
		h := spriteHeight
		px, py := place.X, place.Y
		ps := &PackedSprite{Sprite: s, SheetIndex: len(sheets), Position: image.Pt(px, py)}
		current.Sprites = append(current.Sprites, ps)
		current.W = max(current.W, px+w)
		current.H = max(current.H, py+h)
	}
	if len(current.Sprites) > 0 || len(sheets) == 0 {
		finalizeSheetSize(current)
		sheets = append(sheets, current)
	}
	var out []*PackedSprite
	for _, sh := range sheets {
		out = append(out, sh.Sprites...)
	}
	return out, sheets, nil
}

func packMaxRectsBestBin(sprites []*Sprite, maxSize, padding int, heuristic maxRectsHeuristic) ([]*PackedSprite, []*Sheet, error) {
	var (
		sheets []*Sheet
		bins   []*maxRectsBin
	)

	for _, s := range sprites {
		spriteWidth, spriteHeight, err := validateSpriteSize(s, maxSize, padding)
		if err != nil {
			return nil, nil, err
		}

		bestBinIdx := -1
		bestPlace := maxrectsPlace{}
		for i, bin := range bins {
			place := bin.findPositionWithHeuristic(spriteWidth+padding, spriteHeight+padding, heuristic)
			if !place.ok {
				continue
			}
			if bestBinIdx == -1 ||
				betterMaxRectsPlaceByHeuristic(place, bestPlace, heuristic) ||
				(sameMaxRectsScoreByHeuristic(place, bestPlace, heuristic) && i < bestBinIdx) {
				bestBinIdx = i
				bestPlace = place
			}
		}

		if bestBinIdx == -1 {
			bin := newMaxRectsBin(maxSize, maxSize)
			place := bin.insertWithHeuristic(spriteWidth+padding, spriteHeight+padding, heuristic)
			if !place.ok {
				return nil, nil, fmt.Errorf("failed to place sprite %q (%dx%d) into a new sheet of size %dx%d", s.Name, spriteWidth, spriteHeight, maxSize, maxSize)
			}

			bins = append(bins, bin)
			sheets = append(sheets, &Sheet{})
			bestBinIdx = len(bins) - 1
			bestPlace = place
		} else {
			bins[bestBinIdx].place(bestPlace)
		}

		ps := &PackedSprite{
			Sprite:     s,
			SheetIndex: bestBinIdx,
			Position:   image.Pt(bestPlace.X, bestPlace.Y),
		}
		sheet := sheets[bestBinIdx]
		sheet.Sprites = append(sheet.Sprites, ps)
		sheet.W = max(sheet.W, bestPlace.X+spriteWidth)
		sheet.H = max(sheet.H, bestPlace.Y+spriteHeight)
	}

	if len(sheets) == 0 {
		empty := &Sheet{}
		finalizeSheetSize(empty)
		sheets = append(sheets, empty)
	}
	for _, sh := range sheets {
		finalizeSheetSize(sh)
	}

	var out []*PackedSprite
	for _, sh := range sheets {
		out = append(out, sh.Sprites...)
	}
	return out, sheets, nil
}

func packMaxRectsFillSheet(sprites []*Sprite, maxSize, padding int, heuristic maxRectsHeuristic) ([]*PackedSprite, []*Sheet, error) {
	var sheets []*Sheet
	remaining := make([]*Sprite, len(sprites))
	copy(remaining, sprites)

	for len(remaining) > 0 {
		current := &Sheet{}
		bin := newMaxRectsBin(maxSize, maxSize)

		for {
			bestSpriteIdx := -1
			bestSpriteArea := -1
			bestPlace := maxrectsPlace{}

			for i, s := range remaining {
				spriteWidth, spriteHeight, err := validateSpriteSize(s, maxSize, padding)
				if err != nil {
					return nil, nil, err
				}
				place := bin.findPositionWithHeuristic(spriteWidth+padding, spriteHeight+padding, heuristic)
				if !place.ok {
					continue
				}
				spriteArea := spriteWidth * spriteHeight
				if bestSpriteIdx == -1 ||
					betterMaxRectsPlaceByHeuristic(place, bestPlace, heuristic) ||
					(sameMaxRectsScoreByHeuristic(place, bestPlace, heuristic) && spriteArea > bestSpriteArea) {
					bestSpriteIdx = i
					bestSpriteArea = spriteArea
					bestPlace = place
				}
			}

			if bestSpriteIdx == -1 {
				break
			}

			s := remaining[bestSpriteIdx]
			spriteWidth, spriteHeight := spriteSize(s)
			bin.place(bestPlace)

			ps := &PackedSprite{
				Sprite:     s,
				SheetIndex: len(sheets),
				Position:   image.Pt(bestPlace.X, bestPlace.Y),
			}
			current.Sprites = append(current.Sprites, ps)
			current.W = max(current.W, bestPlace.X+spriteWidth)
			current.H = max(current.H, bestPlace.Y+spriteHeight)

			remaining = append(remaining[:bestSpriteIdx], remaining[bestSpriteIdx+1:]...)
		}

		if len(current.Sprites) == 0 {
			s := remaining[0]
			spriteWidth, spriteHeight := spriteSize(s)
			return nil, nil, fmt.Errorf("failed to place sprite %q (%dx%d) into a new sheet of size %dx%d", s.Name, spriteWidth, spriteHeight, maxSize, maxSize)
		}

		finalizeSheetSize(current)
		sheets = append(sheets, current)
	}

	if len(sheets) == 0 {
		empty := &Sheet{}
		finalizeSheetSize(empty)
		sheets = append(sheets, empty)
	}

	var out []*PackedSprite
	for _, sh := range sheets {
		out = append(out, sh.Sprites...)
	}
	return out, sheets, nil
}

func validateSpriteSize(s *Sprite, maxSize, padding int) (int, int, error) {
	spriteWidth, spriteHeight := spriteSize(s)
	if spriteWidth+padding > maxSize || spriteHeight+padding > maxSize {
		return 0, 0, fmt.Errorf("sprite %q (%dx%d) exceeds max sheet size %d with padding %d", s.Name, spriteWidth, spriteHeight, maxSize, padding)
	}
	return spriteWidth, spriteHeight, nil
}

func spriteSize(s *Sprite) (int, int) {
	return s.Trimmed.Bounds().Dx(), s.Trimmed.Bounds().Dy()
}

func sortSprites(sprites []*Sprite, less func(a, b *Sprite) bool) []*Sprite {
	sorted := make([]*Sprite, len(sprites))
	copy(sorted, sprites)
	sort.Slice(sorted, func(i, j int) bool {
		return less(sorted[i], sorted[j])
	})
	return sorted
}

func lessByMaxSideThenArea(a, b *Sprite) bool {
	aw, ah := spriteSize(a)
	bw, bh := spriteSize(b)
	amax, bmax := max(aw, ah), max(bw, bh)
	if amax != bmax {
		return amax > bmax
	}
	return aw*ah > bw*bh
}

func lessByAreaThenMaxSide(a, b *Sprite) bool {
	aw, ah := spriteSize(a)
	bw, bh := spriteSize(b)
	aArea, bArea := aw*ah, bw*bh
	if aArea != bArea {
		return aArea > bArea
	}
	return max(aw, ah) > max(bw, bh)
}

type packingScore struct {
	totalArea    int64
	sheetCount   int
	maxSheetArea int64
}

func scorePacking(sheets []*Sheet) packingScore {
	var totalArea int64
	var maxSheetArea int64
	for _, sh := range sheets {
		area := int64(sh.W) * int64(sh.H)
		totalArea += area
		if area > maxSheetArea {
			maxSheetArea = area
		}
	}
	return packingScore{
		totalArea:    totalArea,
		sheetCount:   len(sheets),
		maxSheetArea: maxSheetArea,
	}
}

func betterPackingScore(a, b packingScore) bool {
	if a.sheetCount != b.sheetCount {
		return a.sheetCount < b.sheetCount
	}
	if a.totalArea != b.totalArea {
		return a.totalArea < b.totalArea
	}
	return a.maxSheetArea < b.maxSheetArea
}

func finalizeSheetSize(sh *Sheet) {
	sh.W = max(sh.W, 1)
	sh.H = max(sh.H, 1)
}
