package starmap

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

type Starmap struct {
	Bounds image.Rectangle
	Stars  []Star
}

func (sm Starmap) AddStar(s ...Star) {
	sm.Stars = append(sm.Stars, s...)
}

func (sm Starmap) ToImage() *image.NRGBA64 {
	img := image.NewNRGBA64(sm.Bounds)
	bounds := sm.Bounds

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := color.RGBA{R: 255, G: 255, B: 255, A: 255}
			if sm.IntersectWithStar(x, y) {
				c = color.RGBA{R: 0, G: 0, B: 0, A: 255}
			}
			img.Set(x, y, c)
		}
	}

	return img
}

func (sm Starmap) IntersectWithStar(x, y int) bool {
	for i := range sm.Stars {
		if sm.Stars[i].IntersectWith(x, y) {
			return true
		}
	}

	return false
}

func (sm Starmap) WriteFile(filename string) error {
	f, _ := os.Create(filename)
	defer f.Close()

	return png.Encode(f, sm.ToImage())
}

func (sm Starmap) GetOverlap(sm2 Starmap) float64 {
	overlap := 0.0
	for i := range sm.Stars {
		for j := range sm2.Stars {
			overlap += sm.Stars[i].GetOverlap(sm2.Stars[j])
		}
	}

	return overlap / float64(len(sm.Stars))
}

func (sm Starmap) Offset(x, y int) Starmap {
	for i := range sm.Stars {
		sm.Stars[i].Point.X += x
		sm.Stars[i].Point.Y += y
	}

	return sm
}

func (sm Starmap) Rotate(deg float64) Starmap {
	/*
	   angle = (angle ) * (Math.PI/180); // Convert to radians

	   var rotatedX = Math.cos(angle) * (point.x - center.x) - Math.sin(angle) * (point.y-center.y) + center.x;

	   var rotatedY = Math.sin(angle) * (point.x - center.x) + Math.cos(angle) * (point.y - center.y) + center.y;
	*/

	angle := float64(deg) * math.Pi / 180.0
	cosAngle := math.Cos(angle)
	sinAngle := math.Sin(angle)

	centerX := float64(sm.Bounds.Max.X) / 2
	centerY := float64(sm.Bounds.Max.Y) / 2

	for i := range sm.Stars {
		// Rotate by the center of the image.
		x := float64(sm.Stars[i].Point.X)
		y := float64(sm.Stars[i].Point.Y)
		sm.Stars[i].Point.X = (int)(math.Round(cosAngle*(x-centerX) - sinAngle*(y-centerY) + centerX))
		sm.Stars[i].Point.Y = (int)(math.Round(sinAngle*(x-centerX) + cosAngle*(y-centerY) + centerY))
	}

	return sm
}
