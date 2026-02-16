package pack

import (
	"image"

	"github.com/evaneliasyoung/phex/internal/phaser"
)

type maxRectsBin struct {
	phaser.Size
	free []image.Rectangle
	used []image.Rectangle
}

func newMaxRectsBin(w, h int) *maxRectsBin {
	return &maxRectsBin{
		Size: phaser.Size{W: w, H: h},
		free: []image.Rectangle{image.Rect(0, 0, w, h)},
		used: []image.Rectangle{},
	}
}

type maxRectsHeuristic int

const (
	maxRectsShortSideFit maxRectsHeuristic = iota
	maxRectsLongSideFit
	maxRectsAreaFit
	maxRectsBottomLeft
	maxRectsContactPoint
)

type maxrectsPlace struct {
	phaser.Rect
	shortSideFit int
	longSideFit  int
	areaFit      int
	bottomSide   int
	contactScore int
	ok           bool
}

func (b *maxRectsBin) insertWithHeuristic(w, h int, heuristic maxRectsHeuristic) maxrectsPlace {
	place := b.findPositionWithHeuristic(w, h, heuristic)
	if !place.ok {
		return place
	}
	b.place(place)
	return place
}

func (b *maxRectsBin) findPositionWithHeuristic(w, h int, heuristic maxRectsHeuristic) maxrectsPlace {
	best := maxrectsPlace{}
	for _, r := range b.free {
		if w > r.Dx() || h > r.Dy() {
			continue
		}

		x, y := r.Min.X, r.Min.Y
		candidate := maxrectsPlace{
			Rect:         phaser.Rect{X: x, Y: y, W: w, H: h},
			shortSideFit: min(r.Dx()-w, r.Dy()-h),
			longSideFit:  max(r.Dx()-w, r.Dy()-h),
			areaFit:      r.Dx()*r.Dy() - w*h,
			bottomSide:   y + h,
			contactScore: b.contactPointScore(x, y, w, h),
			ok:           true,
		}
		if betterMaxRectsPlaceByHeuristic(candidate, best, heuristic) {
			best = candidate
		}
	}
	return best
}

func (b *maxRectsBin) place(p maxrectsPlace) {
	if !p.ok {
		return
	}
	used := image.Rect(p.X, p.Y, p.X+p.W, p.Y+p.H)
	b.splitFreeRects(used)
	b.prune()
	b.used = append(b.used, used)
}

func betterMaxRectsPlaceByHeuristic(a, b maxrectsPlace, heuristic maxRectsHeuristic) bool {
	if !a.ok {
		return false
	}
	if !b.ok {
		return true
	}

	switch heuristic {
	case maxRectsShortSideFit:
		if a.shortSideFit != b.shortSideFit {
			return a.shortSideFit < b.shortSideFit
		}
		if a.longSideFit != b.longSideFit {
			return a.longSideFit < b.longSideFit
		}
		if a.areaFit != b.areaFit {
			return a.areaFit < b.areaFit
		}
	case maxRectsLongSideFit:
		if a.longSideFit != b.longSideFit {
			return a.longSideFit < b.longSideFit
		}
		if a.shortSideFit != b.shortSideFit {
			return a.shortSideFit < b.shortSideFit
		}
		if a.areaFit != b.areaFit {
			return a.areaFit < b.areaFit
		}
	case maxRectsAreaFit:
		if a.areaFit != b.areaFit {
			return a.areaFit < b.areaFit
		}
		if a.shortSideFit != b.shortSideFit {
			return a.shortSideFit < b.shortSideFit
		}
		if a.longSideFit != b.longSideFit {
			return a.longSideFit < b.longSideFit
		}
	case maxRectsBottomLeft:
		if a.bottomSide != b.bottomSide {
			return a.bottomSide < b.bottomSide
		}
		if a.X != b.X {
			return a.X < b.X
		}
		if a.shortSideFit != b.shortSideFit {
			return a.shortSideFit < b.shortSideFit
		}
	case maxRectsContactPoint:
		if a.contactScore != b.contactScore {
			return a.contactScore > b.contactScore
		}
		if a.areaFit != b.areaFit {
			return a.areaFit < b.areaFit
		}
		if a.shortSideFit != b.shortSideFit {
			return a.shortSideFit < b.shortSideFit
		}
	}

	if a.Y != b.Y {
		return a.Y < b.Y
	}
	return a.X < b.X
}

func sameMaxRectsScoreByHeuristic(a, b maxrectsPlace, heuristic maxRectsHeuristic) bool {
	if !a.ok || !b.ok {
		return false
	}

	switch heuristic {
	case maxRectsShortSideFit:
		return a.shortSideFit == b.shortSideFit &&
			a.longSideFit == b.longSideFit &&
			a.areaFit == b.areaFit
	case maxRectsLongSideFit:
		return a.longSideFit == b.longSideFit &&
			a.shortSideFit == b.shortSideFit &&
			a.areaFit == b.areaFit
	case maxRectsAreaFit:
		return a.areaFit == b.areaFit &&
			a.shortSideFit == b.shortSideFit &&
			a.longSideFit == b.longSideFit
	case maxRectsBottomLeft:
		return a.bottomSide == b.bottomSide &&
			a.X == b.X &&
			a.shortSideFit == b.shortSideFit
	case maxRectsContactPoint:
		return a.contactScore == b.contactScore &&
			a.areaFit == b.areaFit &&
			a.shortSideFit == b.shortSideFit
	default:
		return a.shortSideFit == b.shortSideFit &&
			a.longSideFit == b.longSideFit &&
			a.areaFit == b.areaFit
	}
}

func (b *maxRectsBin) contactPointScore(x, y, w, h int) int {
	score := 0
	if x == 0 || x+w == b.W {
		score += h
	}
	if y == 0 || y+h == b.H {
		score += w
	}

	for _, used := range b.used {
		if used.Min.X == x+w || used.Max.X == x {
			score += commonIntervalLength(used.Min.Y, used.Max.Y, y, y+h)
		}
		if used.Min.Y == y+h || used.Max.Y == y {
			score += commonIntervalLength(used.Min.X, used.Max.X, x, x+w)
		}
	}
	return score
}

func commonIntervalLength(a0, a1, b0, b1 int) int {
	start := max(a0, b0)
	end := min(a1, b1)
	if end <= start {
		return 0
	}
	return end - start
}

func (b *maxRectsBin) splitFreeRects(used image.Rectangle) {
	var newFree []image.Rectangle
	for _, r := range b.free {
		if !r.Overlaps(used) {
			newFree = append(newFree, r)
			continue
		}

		// Canonical MaxRects split: keep maximal candidates, even if they overlap.
		if used.Min.X < r.Max.X && used.Max.X > r.Min.X {
			if used.Min.Y > r.Min.Y && used.Min.Y < r.Max.Y {
				newFree = append(newFree, image.Rect(r.Min.X, r.Min.Y, r.Max.X, used.Min.Y))
			}
			if used.Max.Y < r.Max.Y {
				newFree = append(newFree, image.Rect(r.Min.X, used.Max.Y, r.Max.X, r.Max.Y))
			}
		}

		if used.Min.Y < r.Max.Y && used.Max.Y > r.Min.Y {
			if used.Min.X > r.Min.X && used.Min.X < r.Max.X {
				newFree = append(newFree, image.Rect(r.Min.X, r.Min.Y, used.Min.X, r.Max.Y))
			}
			if used.Max.X < r.Max.X {
				newFree = append(newFree, image.Rect(used.Max.X, r.Min.Y, r.Max.X, r.Max.Y))
			}
		}
	}
	b.free = newFree
}

func (b *maxRectsBin) prune() {
	for i := 0; i < len(b.free); i++ {
		a := b.free[i]
		for j := i + 1; j < len(b.free); j++ {
			bRect := b.free[j]
			if contains(a, bRect) {
				b.free = append(b.free[:j], b.free[j+1:]...)
				j--
				continue
			}
			if contains(bRect, a) {
				b.free = append(b.free[:i], b.free[i+1:]...)
				i--
				break
			}
		}
	}
}
