package pack

import (
	"image"

	"github.com/evaneliasyoung/phex/internal/phaser"
)

type maxRectsBin struct {
	phaser.Size
	free []image.Rectangle
}

func newMaxRectsBin(w, h int) *maxRectsBin {
	return &maxRectsBin{Size: phaser.Size{W: w, H: h}, free: []image.Rectangle{image.Rect(0, 0, w, h)}}
}

type maxrectsPlace struct {
	phaser.Rect
	rot bool
	ok  bool
}

func (b *maxRectsBin) insert(w, h int) maxrectsPlace {
	best := maxrectsPlace{ok: false}
	bestSS, bestLS := 1<<30, 1<<30
	for _, r := range b.free {
		if w <= r.Dx() && h <= r.Dy() {
			ss := min(r.Dx()-w, r.Dy()-h)
			ls := max(r.Dx()-w, r.Dy()-h)
			if ss < bestSS || (ss == bestSS && ls < bestLS) {
				best = maxrectsPlace{Rect: phaser.Rect{X: r.Min.X, Y: r.Min.Y, W: w, H: h}, rot: false, ok: true}
				bestSS, bestLS = ss, ls
			}
		}
	}
	if !best.ok {
		return best
	}
	b.splitFreeRects(image.Rect(best.X, best.Y, best.X+best.W, best.Y+best.H))
	b.prune()
	return best
}

func (b *maxRectsBin) splitFreeRects(used image.Rectangle) {
	var newFree []image.Rectangle
	for _, r := range b.free {
		if !r.Overlaps(used) {
			newFree = append(newFree, r)
			continue
		}
		if used.Min.X > r.Min.X {
			newFree = append(newFree, image.Rect(r.Min.X, r.Min.Y, used.Min.X, r.Max.Y))
		}
		if used.Max.X < r.Max.X {
			newFree = append(newFree, image.Rect(used.Max.X, r.Min.Y, r.Max.X, r.Max.Y))
		}
		if used.Min.Y > r.Min.Y {
			newFree = append(newFree, image.Rect(max(r.Min.X, used.Min.X), r.Min.Y, min(r.Max.X, used.Max.X), used.Min.Y))
		}
		if used.Max.Y < r.Max.Y {
			newFree = append(newFree, image.Rect(max(r.Min.X, used.Min.X), used.Max.Y, min(r.Max.X, used.Max.X), r.Max.Y))
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
