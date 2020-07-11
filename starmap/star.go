package starmap

import (
	"math"
)

type Star struct {
	X    float64
	Y    float64
	Size float64
}

func (s Star) IntersectWith(x, y float64) bool {
	xDist := math.Abs(s.X - x)
	yDist := math.Abs(s.Y - y)

	return math.Sqrt(xDist*xDist+yDist*yDist)-s.Size < 0
}

func (s Star) GetOverlap(s2 Star) float64 {
	if s.X == s2.X && s.Y == s2.Y {
		return 1.0
	}
	return 0.0
}
