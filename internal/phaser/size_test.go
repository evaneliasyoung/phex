package phaser

import "testing"

func TestSizeMin(t *testing.T) {
	s := Size{W: 2048, H: 2048}

	ptMin := s.Min()

	if ptMin.X != 0 {
		t.Fatalf("expected X 0, got %v", ptMin.X)
	}

	if ptMin.Y != 0 {
		t.Fatalf("expected Y 0, got %v", ptMin.Y)
	}
}

func TestSizeMax(t *testing.T) {
	s := Size{W: 2048, H: 2048}

	ptMax := s.Max()

	if ptMax.X != 2048 {
		t.Fatalf("expected X 2048, got %v", ptMax.X)
	}

	if ptMax.Y != 2048 {
		t.Fatalf("expected Y 2048, got %v", ptMax.Y)
	}
}

func TestSizeRect(t *testing.T) {
	s := Size{W: 2048, H: 2048}

	imageRect := s.Rect()

	if imageRect.Min != s.Min() {
		t.Fatalf("expected min %v, got %v", s.Min(), imageRect.Min)
	}

	if imageRect.Max != s.Max() {
		t.Fatalf("expected min %v, got %v", s.Max(), imageRect.Max)
	}
}
