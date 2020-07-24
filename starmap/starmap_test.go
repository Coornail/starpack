package starmap

import (
	"image"
	"image/png"
	"os"
	"testing"

	"golang.org/x/exp/errors/fmt"
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
	if overlap <= 0 {
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

	sm2 := sm1.Copy()
	sm1.Stars[0].X = 511
	sm1.Stars[0].Y = 367

	sm2 = sm2.Offset(-1, -1)
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

func TestOffsetPartial(t *testing.T) {
	s1 := Star{X: 5.5, Y: 4.125, Size: 2.8284271247461903}
	s2 := Star{X: 4.142857142857, Y: 1.857142857143, Size: 2.6457513110645907}

	m1 := Starmap{
		Bounds: image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 20, Y: 20}},
		Stars:  Stars{s1},
	}

	m2 := Starmap{
		Bounds: image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 20, Y: 20}},
		Stars:  Stars{s2},
	}

	sm := Starmaps{m1, m2}

	beforePixels := sm.CorrectPixels()
	offset := m1.FindOffset(m2)

	m2 = m2.Offset(float64(offset.X), float64(offset.Y))
	sm = Starmaps{m1, m2}
	afterPixels := sm.CorrectPixels()

	if beforePixels >= afterPixels {
		t.Errorf("Alignment should improve correct pixels")
	}
}

func TestRotation(t *testing.T) {
	s1 := Star{X: 10, Y: 10, Size: 3}
	s2 := Star{X: 10, Y: 5, Size: 3}
	s3 := Star{X: 5, Y: 5, Size: 3}

	m1 := Starmap{
		Bounds: image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 20, Y: 20}},
		Stars:  Stars{s1, s2},
	}

	m2 := Starmap{
		Bounds: image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 20, Y: 20}},
		Stars:  Stars{s1, s3},
	}

	offset := m1.FindOffset(m2)
	fmt.Printf("%#v\n", offset)
	m2 = m2.Rotate(offset.Rotation)

	overlap := m1.GetOverlap(m2)
	fmt.Printf("%#v\n", overlap)

	f, _ := os.Create("m1.png")
	defer f.Close()
	png.Encode(f, m1.ToImage())

	f2, _ := os.Create("m2.png")
	defer f2.Close()
	png.Encode(f2, m2.ToImage())
}
