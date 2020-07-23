package starmap

import (
	"fmt"
	"image"
	"testing"
)

func TestNoOverlap(t *testing.T) {
	sm := Starmap{
		Bounds: image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 1024, Y: 768}},
		Stars: []Star{{
			X:    1,
			Y:    1,
			Size: 25.0,
		}},
	}

	sm2 := Starmap{
		Bounds: image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 1024, Y: 768}},
		Stars: []Star{{
			X:    1024,
			Y:    768,
			Size: 25.0,
		}},
	}

	overlap := sm.GetOverlap(sm2)
	if overlap > 0.0 {
		t.Errorf("Star maps should not overlap")
	}
}

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

func TestPartialOverlap(t *testing.T) {
	sm := Starmap{
		Bounds: image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 1024, Y: 768}},
		Stars: []Star{
			{
				X:    512,
				Y:    368,
				Size: 25.0,
			},
		},
	}

	sm2 := Starmap{
		Bounds: image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 1024, Y: 768}},
		Stars: []Star{
			{
				X:    511,
				Y:    367,
				Size: 25.0,
			},
		},
	}

	overlap := sm.GetOverlap(sm2)
	if !(overlap > 0) {
		t.Errorf("Stars should overlap")
	}
}

func TestInsideOverlap(t *testing.T) {
	sm := Starmap{
		Bounds: image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 1024, Y: 768}},
		Stars: []Star{
			{
				X:    512,
				Y:    368,
				Size: 25.0,
			},
		},
	}

	sm2 := Starmap{
		Bounds: image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 1024, Y: 768}},
		Stars: []Star{
			{
				X:    512,
				Y:    368,
				Size: 12.0,
			},
		},
	}

	overlap := sm.GetOverlap(sm2)
	if !(overlap > 0 && overlap < 1.0) {
		t.Errorf("Stars should overlap")
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

func BenchmarkOverlap(b *testing.B) {
	s1 := Star{
		X:    512,
		Y:    368,
		Size: 25.0,
	}

	for n := 0; n <= b.N; n++ {
		s2 := Star{
			X:    512,
			Y:    float64(368) + float64(n),
			Size: 25.0,
		}

		s1.GetOverlap(s2)
	}

}

func TestTest(t *testing.T) {
	s1 := Star{X: 1505.5, Y: 1584.125, Size: 2.8284271247461903}
	s2 := Star{X: 1504.142857142857, Y: 1581.857142857143, Size: 2.6457513110645907}

	s3 := s2.Offset(1, 3)

	overlap := s1.GetOverlap(s3)
	fmt.Printf("%#v\n", overlap)

}
