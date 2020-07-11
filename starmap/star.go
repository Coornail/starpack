package starmap

import (
	"image"
	"math"
)

type Star struct {
	Point image.Point
	Size  float64
}

func (s Star) IntersectWith(x, y int) bool {
	xDist := math.Abs(float64(s.Point.X) - float64(x))
	yDist := math.Abs(float64(s.Point.Y) - float64(y))

	return math.Sqrt(xDist*xDist+yDist*yDist)-s.Size < 0
}

func (s Star) GetOverlap(s2 Star) float64 {
	if s.Point.X == s2.Point.X && s.Point.Y == s2.Point.Y {
		return 1.0
	}
	return 0.0
}
