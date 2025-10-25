package phaser

import "image"

type Size struct {
	W int `json:"w"`
	H int `json:"h"`
}

func (sz Size) Min() image.Point {
	return image.Point{0, 0}
}

func (sz Size) Max() image.Point {
	return image.Point{sz.W, sz.H}
}

func (sz Size) Rect() image.Rectangle {
	return image.Rectangle{sz.Min(), sz.Max()}
}
