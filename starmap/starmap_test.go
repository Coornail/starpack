package starmap

import (
	"image"
	"testing"
)

func TestFullOverlap(t *testing.T) {
	sm := Starmap{
		Bounds: image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 1024, Y: 768}},
		Stars: []Star{{
			X:    512,
			Y:    368,
			Size: 25.0,
		}},
	}

	overlap := sm.GetOverlap(sm)
	if overlap != 1.0 {
		t.Errorf("Star overlap with itself is not 1.0")
	}
}

func TestOffset(t *testing.T) {
	sm1 := Starmap{
		Bounds: image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 1024, Y: 768}},
		Stars: []Star{{
			X:    512,
			Y:    368,
			Size: 1.0,
		}},
	}

	sm2 := sm1
	sm1.Stars[0].X = 511
	sm1.Stars[0].Y = 367

	sm2 = sm2.Offset(1, 1)
	overlap := sm1.GetOverlap(sm2)
	if overlap != 1.0 {
		t.Errorf("Star offset failed")
	}
}
