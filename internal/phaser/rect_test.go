package phaser

import "testing"

func TestRectMin(t *testing.T) {
	r := Rect{X: 0, Y: 0, W: 2048, H: 2048}

	ptMin := r.Min()

	if ptMin.X != 0 {
		t.Fatalf("expected X 0, got %v", ptMin.X)
	}

	if ptMin.Y != 0 {
		t.Fatalf("expected Y 0, got %v", ptMin.Y)
	}
}

func TestRectMax(t *testing.T) {
	r := Rect{X: 0, Y: 0, W: 2048, H: 2048}

	ptMax := r.Max()

	if ptMax.X != 2048 {
		t.Fatalf("expected X 2048, got %v", ptMax.X)
	}

	if ptMax.Y != 2048 {
		t.Fatalf("expected Y 2048, got %v", ptMax.Y)
	}
}

func TestRectRect(t *testing.T) {
	r := Rect{X: 0, Y: 0, W: 2048, H: 2048}

	imageRect := r.Rect()

	if imageRect.Min != r.Min() {
		t.Fatalf("expected min %v, got %v", r.Min(), imageRect.Min)
	}

	if imageRect.Max != r.Max() {
		t.Fatalf("expected min %v, got %v", r.Max(), imageRect.Max)
	}
}
