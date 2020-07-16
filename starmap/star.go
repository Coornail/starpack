package starmap

import (
	"math"
)

type Star struct {
	X    float64
	Y    float64
	Size float64
}

func (s Star) Copy() Star {
	return Star{
		X:    s.X,
		Y:    s.Y,
		Size: s.Size,
	}
}

func (s Star) IntersectWith(x, y float64) bool {
	xDist := math.Abs(s.X - x)
	yDist := math.Abs(s.Y - y)

	return math.Sqrt(xDist*xDist+yDist*yDist)-s.Size < 0
}

func (s Star) GetOverlap(s2 Star) float64 {
	dx := math.Abs(s2.X - s.X)
	dy := math.Abs(s2.Y - s.Y)
	d := math.Sqrt(dx*dx + dy*dy)

	// Full overlap.
	if d == 0 {
		return 1.0
	}

	// No overlap.
	if d > s.Size+s2.Size {
		return 0.0
	}

	// Partial overlap.
	return (s.Size + s2.Size) / d
}

func (s Star) IsNeighbor(s2 Star) bool {
	sCopy := s.Copy()
	sCopy.Size++

	return sCopy.GetOverlap(s2) > 0
}
