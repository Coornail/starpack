package starmap

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

const minimumStarSize = 1

type Stars []Star

// IsCloseTo checks if a star is neighbouring a cluster.
func (stars Stars) IsCloseTo(s Star) bool {
	for i := range stars {
		if stars[i].IsNeighbor(s) {
			return true
		}
	}
	return false
}

func (stars Stars) Center() Star {
	centerX := 0.0
	centerY := 0.0

	starLength := float64(len(stars))

	for i := range stars {
		centerX += stars[i].X
		centerY += stars[i].Y
	}

	centerX /= starLength
	centerY /= starLength

	return Star{
		X:    centerX,
		Y:    centerY,
		Size: math.Sqrt(starLength),
	}
}

type Starmap struct {
	Bounds image.Rectangle
	Stars  Stars
}

func (sm Starmap) Copy() Starmap {
	stars := make([]Star, len(sm.Stars))
	copy(stars, sm.Stars)

	return Starmap{
		Bounds: sm.Bounds,
		Stars:  stars,
	}
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
			if sm.IntersectWithStar(float64(x), float64(y)) {
				c = color.RGBA{R: 0, G: 0, B: 0, A: 255}
			}
			img.Set(x, y, c)
		}
	}

	return img
}

func (sm Starmap) IntersectWithStar(x, y float64) bool {
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

func (sm Starmap) Offset(x, y float64) Starmap {
	output := sm.Copy()

	for i := range output.Stars {
		output.Stars[i].X += x
		output.Stars[i].Y += y
	}

	return output
}

func (sm Starmap) Rotate(deg float64) Starmap {
	// https://www.gamefromscratch.com/post/2012/11/24/GameDev-math-recipes-Rotating-one-point-around-another-point.aspx
	angle := float64(deg) * math.Pi / 180.0
	cosAngle := math.Cos(angle)
	sinAngle := math.Sin(angle)

	centerX := float64(sm.Bounds.Max.X) / 2.0
	centerY := float64(sm.Bounds.Max.Y) / 2.0

	for i := range sm.Stars {
		// Rotate by the center of the image.
		x := float64(sm.Stars[i].X)
		y := float64(sm.Stars[i].Y)
		sm.Stars[i].X = cosAngle*(x-centerX) - sinAngle*(y-centerY) + centerX
		sm.Stars[i].Y = sinAngle*(x-centerX) + cosAngle*(y-centerY) + centerY
	}

	return sm
}

func (sm Starmap) FindOffset(sm2 Starmap) (OffsetConfig, float64) {
	max := 1000
	m2 := max * max

	var xMotion, yMotion int
	var dx, dy = 0, -1 // Direction.
	var bestX, bestY int
	var bestRotation int
	maxCorrectPixels := -1.0

	for i := 0; i < m2; i++ {
		if (-max/2 < xMotion && xMotion <= max/2) && (-max/2 < yMotion && yMotion <= max/2) {
			for rotation := -10; rotation <= 10; rotation += 1 {
				correctPixels := Starmaps{sm, sm2.Offset(float64(xMotion), float64(yMotion)).Rotate(float64(rotation))}.CorrectPixels()
				// fmt.Printf("X: %d\t Y: %d\t Rotation: %d\t %f\n", xMotion, yMotion, rotation, correctPixels)
				if maxCorrectPixels < correctPixels {
					maxCorrectPixels = correctPixels
					bestX = xMotion
					bestY = yMotion
					bestRotation = rotation
				}
			}
		}

		// Change direction.
		if xMotion == yMotion || (xMotion < 0 && xMotion == -yMotion) || (xMotion > 0 && xMotion == 1-yMotion) {
			dx, dy = -dy, dx
		}
		xMotion, yMotion = xMotion+dx, yMotion+dy
	}

	return OffsetConfig{X: bestX, Y: bestY, Rotation: float64(bestRotation)}, maxCorrectPixels

}

// Compress several stars into appropriate bigger stars.
// Find neighboring stars and add them together.
func (sm Starmap) Compress() Starmap {
	var clusters []Stars
	var foundCluster bool

	for i := range sm.Stars {
		foundCluster = false

		for j := range clusters {
			if !foundCluster && clusters[j].IsCloseTo(sm.Stars[i]) {
				clusters[j] = append(clusters[j], sm.Stars[i])
				foundCluster = true
			}
		}

		if !foundCluster {
			clusters = append(clusters, Stars{sm.Stars[i]})
		}
	}

	var starmap Starmap
	starmap.Bounds = sm.Bounds
	for i := range clusters {
		starmap.Stars = append(starmap.Stars, clusters[i].Center())
	}

	return starmap
}

type Starmaps []Starmap

func (sms Starmaps) IsOverlap(x, y int) bool {
	reference := sms[0]
	referenceHit := reference.IntersectWithStar(float64(x), float64(y))

	for i := 1; i < len(sms); i++ {
		targetHit := sms[i].IntersectWithStar(float64(x), float64(y))

		if referenceHit != targetHit {
			return false
		}
	}

	return true
}

func (sms Starmaps) VisualizeDifference() *image.NRGBA64 {
	img := image.NewNRGBA64(sms[0].Bounds)
	bounds := sms[0].Bounds

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := color.RGBA{R: 255, G: 255, B: 255, A: 255}
			if sms.IsOverlap(x, y) {
				c = color.RGBA{R: 0, G: 255, B: 0, A: 255}
			} else {
				c = color.RGBA{R: 255, G: 0, B: 0, A: 255}
			}

			img.Set(x, y, c)
		}
	}

	return img
}

func (sms Starmaps) CorrectPixels() float64 {
	ref := sms[0]
	target := sms[1]
	highestOverlap := 0.0

	for i := range ref.Stars {
		for j := range target.Stars {
			overlap := ref.Stars[i].GetOverlap(target.Stars[j])
			if overlap > highestOverlap {
				highestOverlap = overlap
			}
		}
	}

	return highestOverlap
}
