package phaser

import "image"

type Rect struct {
	W int `json:"w"`
	H int `json:"h"`
	X int `json:"x"`
	Y int `json:"y"`
}

func (r Rect) Min() image.Point {
	return image.Point{r.X, r.Y}
}

func (r Rect) Max() image.Point {
	return image.Point{r.X + r.W, r.Y + r.H}
}

func (r Rect) Rect() image.Rectangle {
	return image.Rectangle{r.Min(), r.Max()}
}
